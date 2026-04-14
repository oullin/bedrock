package support

import (
	"testing"
)

// Port of Framework\Tests\Support\LotteryTest::it_can_always_win
func TestLotteryAlwaysWin(t *testing.T) {
	t.Parallel()

	won := false
	NewLottery(1, 100).
		Always().
		Winner(func(...any) any { won = true; return nil }).
		Run()

	if !won {
		t.Error("Always() should always win")
	}
}

// Port of Framework\Tests\Support\LotteryTest::it_can_always_lose
func TestLotteryAlwaysLose(t *testing.T) {
	t.Parallel()

	lost := false
	NewLottery(99, 100).
		Never().
		Loser(func(...any) any { lost = true; return nil }).
		Run()

	if !lost {
		t.Error("Never() should always lose")
	}
}

// Port of Framework\Tests\Support\LotteryTest::it_can_reset_force
func TestLotteryResetForce(t *testing.T) {
	t.Parallel()

	l := NewLottery(1, 1) // 100% odds
	l.ForceWin()
	l.ResetForce()

	// With 100% odds, should still win after reset
	if !l.Choose() {
		t.Error("after reset, 1/1 lottery should win")
	}
}

// Port of Framework\Tests\Support\LotteryTest::it_can_force_win
func TestLotteryForceWin(t *testing.T) {
	t.Parallel()

	l := NewLottery(0, 100) // 0% odds
	l.ForceWin()

	if !l.Choose() {
		t.Error("ForceWin() should override zero odds")
	}
}

// Port of Framework\Tests\Support\LotteryTest::it_can_force_lose
func TestLotteryForceLose(t *testing.T) {
	t.Parallel()

	l := NewLottery(100, 100) // 100% odds
	l.ForceLose()

	if l.Choose() {
		t.Error("ForceLose() should override full odds")
	}
}

// Port of Framework\Tests\Support\LotteryTest::it_runs_the_winner_callback
func TestLotteryWinnerCallback(t *testing.T) {
	t.Parallel()

	var result string
	NewLottery(1, 1).
		Winner(func(args ...any) any {
			result = "won"
			return nil
		}).
		Run()

	if result != "won" {
		t.Errorf("winner callback not invoked: %q", result)
	}
}

// Port of Framework\Tests\Support\LotteryTest::it_runs_the_loser_callback
func TestLotteryLoserCallback(t *testing.T) {
	t.Parallel()

	var result string
	NewLottery(0, 1).
		Loser(func(args ...any) any {
			result = "lost"
			return nil
		}).
		Run()

	if result != "lost" {
		t.Errorf("loser callback not invoked: %q", result)
	}
}

// Port of Framework\Tests\Support\LotteryTest::it_passes_arguments_to_callback
func TestLotteryPassesArgs(t *testing.T) {
	t.Parallel()

	var received []any
	NewLottery(1, 1).
		Winner(func(args ...any) any {
			received = args
			return nil
		}).
		Run("hello", 42)

	if len(received) != 2 || received[0] != "hello" || received[1] != 42 {
		t.Errorf("unexpected args: %v", received)
	}
}

// Port of Framework\Tests\Support\LotteryTest::it_can_use_a_fixed_sequence
func TestLotteryFixedSequence(t *testing.T) {
	t.Parallel()

	l := NewLottery(1, 2)
	fixed := l.Fix([]bool{true, false, true})

	results := []bool{
		fixed.Run(),
		fixed.Run(),
		fixed.Run(),
	}

	if results[0] != true || results[1] != false || results[2] != true {
		t.Errorf("unexpected sequence: %v", results)
	}
}

// Port of Framework\Tests\Support\LotteryTest::it_returns_false_when_sequence_is_exhausted
func TestLotterySequenceExhausted(t *testing.T) {
	t.Parallel()

	fixed := NewLottery(1, 1).Fix([]bool{true})
	fixed.Run() // consume the only item

	// Beyond the sequence: returns false (default)
	if fixed.Run() {
		t.Error("exhausted sequence should return false")
	}
}

// Port of Framework\Tests\Support\LotteryTest::odds_alias
func TestLotteryOdds(t *testing.T) {
	t.Parallel()

	l := LotteryOdds(1, 1)
	l.ForceWin()
	if !l.Choose() {
		t.Error("LotteryOdds alias should work")
	}
}
