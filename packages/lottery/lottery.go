package lottery

import (
	"errors"
	"math/rand"
	"sync"
)

// Lottery provides probabilistic execution with configurable odds.
// It mirrors the executable parts of Illuminate\Support\Lottery.
type Lottery struct {
	chances float64
	outOf   *int
	winner  func(...any) any
	loser   func(...any) any
	forced  *bool // nil = use random; true = always win; false = always lose
}

// fixedLottery wraps a Lottery with a predetermined sequence of outcomes.
type fixedLottery struct {
	*Lottery
	sequence []bool
	pos      int
}

var (
	resultFactoryMu sync.RWMutex
	resultFactory   func(chances float64, outOf *int) bool
)

// NewLottery creates a new Lottery with the given chance and optional outOf.
//
// The executable parity here follows Laravel's constructor behavior:
//   - one argument means a floating-point chance in the range [0, 1]
//   - two arguments mean "chance out of outOf"
//
// Go cannot surface the PHP exception type, so invalid arguments panic with a
// clear message.
func NewLottery(chances float64, outOf ...int) *Lottery {
	if len(outOf) > 0 && outOf[0] < 1 {
		panic(errors.New("outOf must be at least 1"))
	}

	if len(outOf) == 0 && chances > 1 {
		panic(errors.New("Float must not be greater than 1."))
	}

	var denom *int

	if len(outOf) > 0 {
		value := outOf[0]
		denom = &value
	}

	return &Lottery{
		chances: chances,
		outOf:   denom,
	}
}

// LotteryOdds matches Laravel's Lottery::odds() helper.
func LotteryOdds(chances float64, outOf ...int) *Lottery {
	return NewLottery(chances, outOf...)
}

// Winner sets the callback to invoke when the lottery wins.
func (l *Lottery) Winner(fn func(...any) any) *Lottery {
	l.winner = fn

	return l
}

// Loser sets the callback to invoke when the lottery loses.
func (l *Lottery) Loser(fn func(...any) any) *Lottery {
	l.loser = fn

	return l
}

// Run executes the lottery and returns the callback result.
func (l *Lottery) Run(args ...any) any {
	return l.runWith(l.determine, args...)
}

// Choose evaluates the lottery once, or multiple times when times is provided.
// Without callbacks the result is a boolean. With callbacks it returns the
// callback return value(s), matching Laravel's executable behavior.
func (l *Lottery) Choose(times ...int) any {
	if len(times) == 0 {
		return l.runWith(l.determine)
	}

	if times[0] <= 0 {
		return []any{}
	}

	results := make([]any, 0, times[0])

	for i := 0; i < times[0]; i++ {
		results = append(results, l.runWith(l.determine))
	}

	return results
}

func (l *Lottery) determine() bool {
	if l.forced != nil {
		return *l.forced
	}

	return currentResultFactory()(l.chances, l.outOf)
}

func (l *Lottery) runWith(determine func() bool, args ...any) any {
	if determine() {
		if l.winner != nil {
			return l.winner(args...)
		}

		return true
	}

	if l.loser != nil {
		return l.loser(args...)
	}

	return false
}

// Always forces the lottery to always win for this instance.
func (l *Lottery) Always() *Lottery {
	t := true
	l.forced = &t

	return l
}

// Never forces the lottery to always lose for this instance.
func (l *Lottery) Never() *Lottery {
	f := false
	l.forced = &f

	return l
}

// ForceWin is a compatibility alias for Always.
func (l *Lottery) ForceWin() {
	l.Always()
}

// ForceLose is a compatibility alias for Never.
func (l *Lottery) ForceLose() {
	l.Never()
}

// ResetForce restores normal probabilistic behavior for this instance.
func (l *Lottery) ResetForce() {
	l.forced = nil
}

// Fix returns a lottery wrapper that replays the supplied outcome sequence.
// This is a local test helper that keeps the existing package API available.
func (l *Lottery) Fix(sequence []bool) *fixedLottery {
	return &fixedLottery{Lottery: l, sequence: sequence}
}

// Run executes the fixed lottery using the next value in the sequence.
func (fl *fixedLottery) Run(args ...any) any {
	return fl.runWith(fl.determine, args...)
}

// Choose evaluates the fixed lottery once or multiple times.
func (fl *fixedLottery) Choose(times ...int) any {
	if len(times) == 0 {
		return fl.runWith(fl.determine)
	}

	if times[0] <= 0 {
		return []any{}
	}

	results := make([]any, 0, times[0])

	for i := 0; i < times[0]; i++ {
		results = append(results, fl.runWith(fl.determine))
	}

	return results
}

func (fl *fixedLottery) determine() bool {
	if fl.pos < len(fl.sequence) {
		won := fl.sequence[fl.pos]
		fl.pos++

		return won
	}

	return false
}

// DetermineResultsNormally removes any package-wide forced result factory.
func DetermineResultsNormally() {
	DetermineResultNormally()
}

// DetermineResultNormally removes any package-wide forced result factory.
func DetermineResultNormally() {
	setResultFactory(nil)
}

// SetResultFactory replaces the package-wide result factory.
func SetResultFactory(factory func(chances float64, outOf *int) bool) {
	setResultFactory(factory)
}

// AlwaysWin forces package-wide lotteries to win while callback executes.
func AlwaysWin(callback func()) {
	withResultFactory(func(float64, *int) bool { return true }, callback)
}

// AlwaysLose forces package-wide lotteries to lose while callback executes.
func AlwaysLose(callback func()) {
	withResultFactory(func(float64, *int) bool { return false }, callback)
}

// ForceResultWithSequence forces package-wide lotteries to use the supplied
// sparse sequence. Missing items fall back to whenMissing, or normal
// probability when whenMissing is nil.
func ForceResultWithSequence(sequence map[int]bool, whenMissing func(float64, *int) bool) {
	var (
		sequenceMu sync.Mutex
		next       int
	)

	setResultFactory(func(chances float64, outOf *int) bool {
		sequenceMu.Lock()

		defer sequenceMu.Unlock()

		if result, ok := sequence[next]; ok {
			next++

			return result
		}

		if whenMissing != nil {
			return whenMissing(chances, outOf)
		}

		next++

		return normalResultFactory(chances, outOf)
	})
}

func withResultFactory(factory func(chances float64, outOf *int) bool, callback func()) {
	resultFactoryMu.Lock()
	previous := resultFactory
	resultFactory = factory
	resultFactoryMu.Unlock()

	defer func() {
		resultFactoryMu.Lock()
		resultFactory = previous
		resultFactoryMu.Unlock()
	}()

	if callback != nil {
		callback()
	}
}

func setResultFactory(factory func(chances float64, outOf *int) bool) {
	resultFactoryMu.Lock()
	resultFactory = factory
	resultFactoryMu.Unlock()
}

func currentResultFactory() func(chances float64, outOf *int) bool {
	resultFactoryMu.RLock()

	defer resultFactoryMu.RUnlock()

	if resultFactory != nil {
		return resultFactory
	}

	return normalResultFactory
}

func normalResultFactory(chances float64, outOf *int) bool {
	if outOf == nil {
		return rand.Float64() <= chances
	}

	if *outOf <= 0 {
		return false
	}

	return float64(rand.Intn(*outOf)+1) <= chances
}
