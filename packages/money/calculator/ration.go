package calculator

import (
	"math"

	"github.com/bedrock/packages/money/exception"
)

// SafeAdd safely adds two int64 amounts, returning 0 on overflow.
func SafeAdd(a, b Amount) Amount {
	if (b > 0 && a > math.MaxInt64-b) || (b < 0 && a < math.MinInt64-b) {
		return 0
	}

	return a + b
}

// SafeSubtract safely subtracts b from a, returning 0 on overflow.
func SafeSubtract(a, b Amount) Amount {
	if (b > 0 && a < math.MinInt64+b) || (b < 0 && a > math.MaxInt64+b) {
		return 0
	}

	return a - b
}

// Ration safely computes the multiplication of an amount by a ration (weight),
// returning 0 if the operation results in an integer overflow.
//
// It is designed to be used in allocation or pro-rating formulas where
// the intermediate calculation (amount * ration) might exceed int64 limits.
//
// Safety Strategy:
// Instead of multiplying first and checking the result (which would be too late),
// this function algebraically rearranges the formula to compare the inputs
// against safe limits before any dangerous multiplication occurs.
//
// Parameters:
//   - amount: The value to be multiplied.
//   - ration: The multiplier/weight.
//
// Returns:
//   - The result of (amount * ration).
//   - 0 if the calculation overflows int64 bounds.
//   - 0 if ration or amount is 0.
func Ration(amount Amount, ration int64) int64 {
	// 1. Guard Clause: Quick exit for zero values.
	if ration == 0 || amount == 0 {
		return 0
	}

	// 2. The -1 Edge Case
	// In 64-bit signed integers, dividing math.MinInt64 by -1 causes an overflow
	// because the positive range is one less than the negative range.
	// (Abs(MinInt64) = MaxInt64 + 1).
	if ration == -1 {
		if amount == math.MinInt64 {
			return 0 // Overflow: Result would be > MaxInt64
		}

		return -amount
	}

	// 3. Overflow Protection Logic
	// We want to ensure: amount * ration <= MaxInt64
	// Rearranged check: amount <= MaxInt64 / ration
	limit := math.MaxInt64 / ration

	overflows := false

	if ration > 0 {
		// Scenario A: Positive Ration
		// The inequality signs remain standard.
		// 1. Max Check: amount > Max / ration (Positive overflow)
		// 2. Min Check: amount < Min / ration (Negative underflow)
		if amount > limit || amount < math.MinInt64/ration {
			overflows = true
		}
	} else {
		// Scenario B: Negative Ration
		// Dividing an inequality by a negative number FLIPS the sign.
		// 1. Max Check: amount < Max / ration (Double negative becomes positive overflow)
		// 2. Min Check: amount > Min / ration (Positive * Negative becomes negative underflow)
		if amount < limit || amount > math.MinInt64/ration {
			overflows = true
		}
	}

	if overflows {
		return 0
	}

	// 4. Safe Calculation
	// If we reach here, the multiplication is guaranteed to fit in int64.
	return amount * ration
}

// SafeMultiply performs sequential multiplication of int64 values with overflow detection.
// It multiplies the initial value by all provided multipliers in sequence, checking for
// overflow at each step.
//
// Parameters:
//   - initial: The starting value to multiply
//   - multipliers: One or more values to multiply by in sequence
//
// Returns:
//   - The final result after all multiplications
//   - An error if any multiplication would cause an overflow
//
// Example:
//
//	result, err := SafeMultiply(10, 2, 3, 4) // 10 * 2 * 3 * 4 = 240
func SafeMultiply(initial int64, multipliers ...int64) (int64, error) {
	result := initial

	for _, value := range multipliers {
		// Skip overflow checks if either value is zero (result will be zero)
		if value == 0 || result == 0 {
			result *= value

			continue
		}

		// Handle the -1 edge case (MinInt64 * -1 overflows)
		if value == -1 && result == math.MinInt64 {
			return 0, exception.ErrOverflow
		}

		if result == -1 && value == math.MinInt64 {
			return 0, exception.ErrOverflow
		}

		// Check if multiplication would overflow for positive multipliers
		if value > 0 {
			if result > math.MaxInt64/value || result < math.MinInt64/value {
				return 0, exception.ErrOverflow
			}

			result *= value

			continue
		}

		// Check if multiplication would overflow for negative multipliers (excluding -1)
		if value < -1 {
			// When dividing by negative, inequality flips
			if result < math.MaxInt64/value || result > math.MinInt64/value {
				return 0, exception.ErrOverflow
			}
		}

		// For value == -1, we already handled the overflow case above,
		// so the multiplication is safe
		result *= value
	}

	return result, nil
}
