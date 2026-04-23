package featureflags

import "math/rand"

// Lottery samples a boolean feature value. It is a small Go equivalent of
// Upstream's lottery-backed feature values.
type Lottery struct {
	draw func() bool
}

// NewLottery creates a Lottery from a caller-provided draw function.
func NewLottery(draw func() bool) Lottery {
	return Lottery{draw: draw}
}

// LotteryOdds creates a Lottery that wins chances out of total draws.
func LotteryOdds(chances, total int) Lottery {
	return NewLottery(func() bool {
		if total <= 0 || chances <= 0 {
			return false
		}

		if chances >= total {
			return true
		}

		return rand.Intn(total) < chances
	})
}

// FixedLottery creates a deterministic Lottery that returns the supplied
// sequence, then repeats the last value.
func FixedLottery(results ...bool) Lottery {
	index := 0

	return NewLottery(func() bool {
		if len(results) == 0 {
			return false
		}

		if index >= len(results) {
			return results[len(results)-1]
		}

		result := results[index]
		index++

		return result
	})
}

func resolveFeatureValue(value any) any {
	switch v := value.(type) {
	case Lottery:
		return v.draw()
	case *Lottery:
		if v == nil {
			return nil
		}

		return v.draw()
	default:
		return value
	}
}
