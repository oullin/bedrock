package calculator

import "math"

// Amount represents the integer value of a monetary amount.
type Amount = int64

// Calculator provides safe arithmetic operations for monetary amounts.
type Calculator struct{}

// NewCalculator creates a new instance of Calculator.
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Add adds two amounts safely, handling potential overflows.
func (c *Calculator) Add(a, b Amount) Amount {
	return SafeAdd(a, b)
}

// Subtract subtracts the second amount from the first safely, handling potential overflows.
func (c *Calculator) Subtract(a, b Amount) Amount {
	return SafeSubtract(a, b)
}

// Multiply multiplies an amount by a given seed (ration).
func (c *Calculator) Multiply(amount Amount, seed int64) Amount {
	return Ration(amount, seed)
}

// SafeMultiply multiplies an initial amount by a series of multipliers, checking for overflow.
func (c *Calculator) SafeMultiply(initial int64, multipliers ...int64) (int64, error) {
	if c == nil {
		return 0, nil
	}

	return SafeMultiply(initial, multipliers...)
}

// Divide divides an amount by a seed. Returns 0 if the seed is 0.
func (c *Calculator) Divide(amount Amount, seed int64) Amount {
	if c == nil || seed == 0 {
		return 0
	}

	return amount / seed
}

// Modulus returns the remainder of dividing an amount by a seed. Returns 0 if the seed is 0.
func (c *Calculator) Modulus(amount Amount, seed int64) Amount {
	if c == nil || seed == 0 {
		return 0
	}

	return amount % seed
}

// Allocate allocates an amount based on a ration and scale.
func (c *Calculator) Allocate(amount Amount, ration, scale int64) Amount {
	if c == nil || amount == 0 || scale == 0 {
		return 0
	}

	return Ration(amount, ration) / scale
}

// Absolute returns the absolute value of an amount.
func (c *Calculator) Absolute(amount Amount) Amount {
	if c == nil || amount < math.MinInt64 {
		return 0
	}

	if amount < 0 {
		return -amount
	}

	return amount
}

// Negative returns the negative value of an amount.
func (c *Calculator) Negative(amount Amount) Amount {
	if c == nil {
		return 0
	}

	if amount > 0 {
		return -amount
	}

	return amount
}

// Round rounds an amount to a specified exponent (precision).
func (c *Calculator) Round(amount Amount, exponent int) Amount {
	// Guard against nil calculator to prevent panics
	if c == nil {
		return 0
	}

	// Return amount unchanged for edge cases:
	// - Zero amounts don't need rounding
	// - Non-positive exponents would cause invalid rounding (no decimal places or negative precision)
	// - Exponents > 18 exceed int64 safe range for 10^exponent calculation
	if amount == 0 || exponent <= 0 || exponent > 18 {
		return amount
	}

	// Work with absolute value to apply consistent rounding logic regardless of sign
	absolute := c.Absolute(amount)

	// Calculate the rounding unit: 10^exponent determines the precision level
	// e.g. exponent=2 means round to the nearest 100 (cents to dollars for precision 2)
	reminder := int64(math.Pow(10, float64(exponent)))

	// Get the remainder to determine if we need to round up or down
	module := absolute % reminder

	// Implement round-to-nearest with ties rounded down (toward zero):
	// only remainders strictly greater than half the unit trigger a round up
	if module > (reminder / 2) {
		absolute += reminder
	}

	// Truncate to the nearest multiple of the rounding unit by dividing and multiplying
	// This effectively removes all digits below the desired precision level
	absolute = (absolute / reminder) * reminder

	// Reapply the original sign to preserve whether the amount was positive or negative
	if amount < 0 {
		amount = -absolute
	} else {
		amount = absolute
	}

	return amount
}
