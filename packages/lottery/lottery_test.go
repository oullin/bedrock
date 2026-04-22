package lottery

import (
	"fmt"
	"reflect"
	"testing"
)

func mustPanic(t *testing.T, want string, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()

		if recovered == nil {
			t.Fatalf("expected panic %q", want)
		}

		if got := fmt.Sprint(recovered); got != want {
			t.Fatalf("panic = %q, want %q", got, want)
		}
	}()

	fn()
}

// Port of Illuminate\Tests\Support\LotteryTest::testItCanWin.
func TestItCanWin(t *testing.T) {
	won := false

	result := NewLottery(1, 1).
		Winner(func(...any) any {
			won = true

			return nil
		}).
		Choose()

	if !won {
		t.Fatal("winner callback was not invoked")
	}

	if result != nil {
		t.Fatalf("Choose() = %v, want nil from winner callback", result)
	}
}

// Port of Illuminate\Tests\Support\LotteryTest::testItCanLose.
func TestItCanLose(t *testing.T) {
	won := false
	lost := false

	result := NewLottery(0, 1).
		Winner(func(...any) any {
			won = true

			return nil
		}).
		Loser(func(...any) any {
			lost = true

			return nil
		}).
		Choose()

	if won {
		t.Fatal("winner callback should not be invoked")
	}

	if !lost {
		t.Fatal("loser callback was not invoked")
	}

	if result != nil {
		t.Fatalf("Choose() = %v, want nil from loser callback", result)
	}
}

// Port of Illuminate\Tests\Support\LotteryTest::testItCanReturnValues.
func TestItCanReturnValues(t *testing.T) {
	win := NewLottery(1, 1).
		Winner(func(...any) any {
			return "win"
		}).
		Choose()

	if win != "win" {
		t.Fatalf("win = %v, want %q", win, "win")
	}

	lose := NewLottery(0, 1).
		Loser(func(...any) any {
			return "lose"
		}).
		Choose()

	if lose != "lose" {
		t.Fatalf("lose = %v, want %q", lose, "lose")
	}
}

// Port of Illuminate\Tests\Support\LotteryTest::testItCanChooseSeveralTimes.
func TestItCanChooseSeveralTimes(t *testing.T) {
	winResults := NewLottery(1, 1).
		Winner(func(...any) any {
			return "win"
		}).
		Choose(2)

	if !reflect.DeepEqual(winResults, []any{"win", "win"}) {
		t.Fatalf("winResults = %#v, want %#v", winResults, []any{"win", "win"})
	}

	loseResults := NewLottery(0, 1).
		Loser(func(...any) any {
			return "lose"
		}).
		Choose(2)

	if !reflect.DeepEqual(loseResults, []any{"lose", "lose"}) {
		t.Fatalf("loseResults = %#v, want %#v", loseResults, []any{"lose", "lose"})
	}
}

// Port of Illuminate\Tests\Support\LotteryTest::testItCanBePassedAsCallable.
func TestItCanBePassedAsCallable(t *testing.T) {
	result := func(callable func(...any) any) any {
		return callable("winner-chicken", "-dinner")
	}(NewLottery(1, 1).
		Winner(func(args ...any) any {
			return "winner-" + args[0].(string) + args[1].(string)
		}).
		Run)

	if result != "winner-winner-chicken-dinner" {
		t.Fatalf("result = %v, want %q", result, "winner-winner-chicken-dinner")
	}
}

// Port of Illuminate\Tests\Support\LotteryTest::testWithoutSpecifiedClosuresBooleansAreReturned.
func TestWithoutSpecifiedClosuresBooleansAreReturned(t *testing.T) {
	win := NewLottery(1, 1).Choose()

	if win != true {
		t.Fatalf("win = %v, want true", win)
	}

	lose := NewLottery(0, 1).Choose()

	if lose != false {
		t.Fatalf("lose = %v, want false", lose)
	}
}

// Port of Illuminate\Tests\Support\LotteryTest::testItCanForceWinningResultInTests.
func TestItCanForceWinningResultInTests(t *testing.T) {
	t.Cleanup(DetermineResultsNormally)

	var result any

	AlwaysWin(func() {
		result = NewLottery(1, 2).
			Winner(func(...any) any {
				return "winner"
			}).
			Choose(10)
	})

	want := []any{
		"winner", "winner", "winner", "winner", "winner",
		"winner", "winner", "winner", "winner", "winner",
	}

	if !reflect.DeepEqual(result, want) {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

// Port of Illuminate\Tests\Support\LotteryTest::testItCanForceLosingResultInTests.
func TestItCanForceLosingResultInTests(t *testing.T) {
	t.Cleanup(DetermineResultsNormally)

	var result any

	AlwaysLose(func() {
		result = NewLottery(1, 2).
			Loser(func(...any) any {
				return "loser"
			}).
			Choose(10)
	})

	want := []any{
		"loser", "loser", "loser", "loser", "loser",
		"loser", "loser", "loser", "loser", "loser",
	}

	if !reflect.DeepEqual(result, want) {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

// Port of Illuminate\Tests\Support\LotteryTest::testItCanForceTheResultViaSequence.
func TestItCanForceTheResultViaSequence(t *testing.T) {
	t.Cleanup(DetermineResultsNormally)

	sequence := map[int]bool{
		0: true,
		1: false,
		2: true,
		3: false,
		4: true,
		5: false,
		6: true,
		7: false,
		8: true,
		9: false,
	}

	ForceResultWithSequence(sequence, nil)

	result := NewLottery(1, 100).
		Winner(func(...any) any {
			return "winner"
		}).
		Loser(func(...any) any {
			return "loser"
		}).
		Choose(10)

	want := []any{
		"winner", "loser", "winner", "loser", "winner",
		"loser", "winner", "loser", "winner", "loser",
	}

	if !reflect.DeepEqual(result, want) {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

// Port of Illuminate\Tests\Support\LotteryTest::testItCanHandleMissingSequenceItems.
func TestItCanHandleMissingSequenceItems(t *testing.T) {
	t.Cleanup(DetermineResultsNormally)

	ForceResultWithSequence(map[int]bool{
		0: true,
		1: true,
		3: true,
	}, func(float64, *int) bool {
		panic("Missing key in sequence.")
	})

	first := NewLottery(1, 10000).
		Winner(func(...any) any {
			return "winner"
		}).
		Loser(func(...any) any {
			return "loser"
		}).
		Choose()

	if first != "winner" {
		t.Fatalf("first = %v, want %q", first, "winner")
	}

	second := NewLottery(1, 10000).
		Winner(func(...any) any {
			return "winner"
		}).
		Loser(func(...any) any {
			return "loser"
		}).
		Choose()

	if second != "winner" {
		t.Fatalf("second = %v, want %q", second, "winner")
	}

	mustPanic(t, "Missing key in sequence.", func() {
		NewLottery(1, 10000).
			Winner(func(...any) any {
				return "winner"
			}).
			Loser(func(...any) any {
				return "loser"
			}).
			Choose()
	})
}

// Port of Illuminate\Tests\Support\LotteryTest::testItThrowsForFloatsOverOne.
func TestItThrowsForFloatsOverOne(t *testing.T) {
	mustPanic(t, "Float must not be greater than 1.", func() {
		NewLottery(1.1)
	})
}

// Port of Illuminate\Tests\Support\LotteryTest::testItThrowsForOutOfLessThanOne.
func TestItThrowsForOutOfLessThanOne(t *testing.T) {
	mustPanic(t, "outOf must be at least 1", func() {
		NewLottery(1, 0)
	})
}

// Port of Illuminate\Tests\Support\LotteryTest::testItCanWinWithFloat.
func TestItCanWinWithFloat(t *testing.T) {
	wins := false

	result := LotteryOdds(1.0).
		Winner(func(...any) any {
			wins = true

			return nil
		}).
		Choose()

	if !wins {
		t.Fatal("winner callback was not invoked for float odds")
	}

	if result != nil {
		t.Fatalf("Choose() = %v, want nil from winner callback", result)
	}
}

// Port of Illuminate\Tests\Support\LotteryTest::testItCanLoseWithFloat.
func TestItCanLoseWithFloat(t *testing.T) {
	wins := false
	loses := false

	result := LotteryOdds(0.0).
		Winner(func(...any) any {
			wins = true

			return nil
		}).
		Loser(func(...any) any {
			loses = true

			return nil
		}).
		Choose()

	if wins {
		t.Fatal("winner callback should not be invoked for float odds")
	}

	if !loses {
		t.Fatal("loser callback was not invoked for float odds")
	}

	if result != nil {
		t.Fatalf("Choose() = %v, want nil from loser callback", result)
	}
}
