package lottery

import (
	"math/rand"
)

// Lottery provides probabilistic execution with configurable odds.
// Mirrors Framework\Support\Lottery.
type Lottery struct {
	numerator   int
	denominator int
	winner      func(...any) any
	loser       func(...any) any
	forced      *bool // nil = use random; true = always win; false = always lose
}

// NewLottery creates a new Lottery with the given chances (numerator out of denominator).
// Example: NewLottery(1, 100) = 1% chance of winning.
// Mirrors new Lottery($chances, $outOf) / Lottery::odds().

// LotteryOdds is an alias for NewLottery that matches Upstream's Lottery::odds() API.

// Winner sets the callback to invoke when the lottery is won.
// Mirrors Lottery::winner().

// Loser sets the callback to invoke when the lottery is lost.
// Mirrors Lottery::loser().

// Run runs the lottery and invokes the winner or loser callback.
// Extra arguments are forwarded to the callbacks.
// Returns true if the lottery was won.
// Mirrors Lottery::__invoke() / choose().

// Choose determines if the lottery is won without invoking callbacks.

//nolint:gosec

// Always forces the lottery to always be won (for testing).
// Mirrors Lottery::alwaysWin().

// Never forces the lottery to always be lost (for testing).
// Mirrors Lottery::alwaysLose().

// ForceWin forces the lottery to always be won.

// ForceLose forces the lottery to always be lost.

// ResetForce restores normal random behaviour.

// Fix sets a sequence of results (true = win, false = lose).
// Used for deterministic testing.
// Mirrors Lottery::fix().

// fixedLottery wraps a Lottery with a predetermined sequence of outcomes.
type fixedLottery struct {
	*Lottery
	sequence []bool
	pos      int
}

func NewLottery(numerator, denominator int) *Lottery {
	return &Lottery{
		numerator:   numerator,
		denominator: denominator,
	}
}

func LotteryOdds(chances, outOf int) *Lottery {
	return NewLottery(chances, outOf)
}

func (l *Lottery) Winner(fn func(...any) any) *Lottery {
	l.winner = fn

	return l
}

func (l *Lottery) Loser(fn func(...any) any) *Lottery {
	l.loser = fn

	return l
}

func (l *Lottery) Run(args ...any) bool {
	won := l.determine()

	if won && l.winner != nil {
		l.winner(args...)
	} else if !won && l.loser != nil {
		l.loser(args...)
	}

	return won
}

func (l *Lottery) Choose() bool {
	return l.determine()
}

func (l *Lottery) determine() bool {
	if l.forced != nil {
		return *l.forced
	}

	if l.denominator <= 0 {
		return false
	}

	return rand.Intn(l.denominator) < l.numerator
}

func (l *Lottery) Always() *Lottery {
	t := true
	l.forced = &t

	return l
}

func (l *Lottery) Never() *Lottery {
	f := false
	l.forced = &f

	return l
}

func (l *Lottery) ForceWin() {
	l.Always()
}

func (l *Lottery) ForceLose() {
	l.Never()
}

func (l *Lottery) ResetForce() {
	l.forced = nil
}

func (l *Lottery) Fix(sequence []bool) *fixedLottery {
	return &fixedLottery{Lottery: l, sequence: sequence}
}

// Run runs the lottery using the next value in the sequence.
func (fl *fixedLottery) Run(args ...any) bool {
	won := false

	if fl.pos < len(fl.sequence) {
		won = fl.sequence[fl.pos]
		fl.pos++
	}

	if won && fl.winner != nil {
		fl.winner(args...)
	} else if !won && fl.loser != nil {
		fl.loser(args...)
	}

	return won
}
