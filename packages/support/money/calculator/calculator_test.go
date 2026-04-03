package calculator

import (
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	c := NewCalculator()

	if c == nil {
		t.Error("NewCalculator() returned nil")
	}
}

func TestAdd(t *testing.T) {
	c := NewCalculator()
	tests := []struct {
		name string
		a, b Amount
		want Amount
	}{
		{"overflow positive", math.MaxInt64, 1, 0},
		{"overflow negative", math.MinInt64, -1, 0},
		{"boundary safe", math.MaxInt64 - 1, 1, math.MaxInt64},
		{"positive numbers", 100, 50, 150},
		{"negative numbers", -100, -50, -150},
		{"mixed numbers", 100, -50, 50},
		{"zero", 100, 0, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Add(tt.a, tt.b); got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	c := NewCalculator()
	tests := []struct {
		name string
		a, b Amount
		want Amount
	}{
		{"overflow positive", math.MaxInt64, -1, 0},
		{"overflow negative", math.MinInt64, 1, 0},
		{"overflow minint64 subtracted from zero", 0, math.MinInt64, 0},
		{"overflow minint64 subtracted from max", math.MaxInt64, math.MinInt64, 0},
		{"minint64 special case boundary", -1, math.MinInt64, math.MaxInt64},
		{"boundary safe", math.MinInt64 + 1, 1, math.MinInt64},
		{"positive result", 100, 50, 50},
		{"negative result", 50, 100, -50},
		{"negative numbers", -50, -100, 50},
		{"zero", 100, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Subtract(tt.a, tt.b); got != tt.want {
				t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	c := NewCalculator()
	tests := []struct {
		name   string
		amount Amount
		seed   int64
		want   Amount
	}{
		{"overflow positive", math.MaxInt64, 2, 0},
		{"overflow negative", math.MinInt64, 2, 0},
		{"boundary safe", math.MaxInt64 / 2, 2, math.MaxInt64 - 1},
		{"positive", 100, 2, 200},
		{"negative multiplier", 100, -2, -200},
		{"zero", 100, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Multiply(tt.amount, tt.seed); got != tt.want {
				t.Errorf("Multiply(%d, %d) = %d; want %d", tt.amount, tt.seed, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	c := NewCalculator()
	tests := []struct {
		name   string
		amount Amount
		seed   int64
		want   Amount
	}{
		{"exact division", 100, 2, 50},
		{"integer division", 100, 3, 33},
		{"divide by zero", 100, 0, 0},
		{"negative result", 100, -2, -50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Divide(tt.amount, tt.seed); got != tt.want {
				t.Errorf("Divide(%d, %d) = %d; want %d", tt.amount, tt.seed, got, tt.want)
			}
		})
	}
}

func TestModulus(t *testing.T) {
	c := NewCalculator()
	tests := []struct {
		name   string
		amount Amount
		seed   int64
		want   Amount
	}{
		{"no remainder", 100, 2, 0},
		{"remainder", 100, 3, 1},
		{"modulus by zero", 100, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Modulus(tt.amount, tt.seed); got != tt.want {
				t.Errorf("Modulus(%d, %d) = %d; want %d", tt.amount, tt.seed, got, tt.want)
			}
		})
	}
}

func TestAllocate(t *testing.T) {
	c := NewCalculator()
	tests := []struct {
		name          string
		calc          *Calculator
		amount        Amount
		ration, scale int64
		want          Amount
	}{
		{"simple allocation", c, 100, 1, 2, 50},
		{"complex allocation", c, 100, 1, 3, 33},
		{"zero amount", c, 0, 1, 2, 0},
		{"zero splitter", c, 100, 1, 0, 0},
		{"nil calculator", nil, 100, 1, 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.calc.Allocate(tt.amount, tt.ration, tt.scale); got != tt.want {
				t.Errorf("Allocate(%d, %d, %d) = %d; want %d", tt.amount, tt.ration, tt.scale, got, tt.want)
			}
		})
	}
}

func TestAllocateOverflow(t *testing.T) {
	c := NewCalculator()
	tests := []struct {
		name          string
		amount        Amount
		ration, scale int64
		want          Amount
	}{
		// Positive overflow: amount * ration > MaxInt64
		{"positive overflow", math.MaxInt64, 2, 1, 0},
		{"positive overflow edge", math.MaxInt64 / 2, 3, 1, 0},

		// Negative overflow: amount * ration < MinInt64
		{"negative overflow", math.MinInt64, 2, 1, 0},
		{"negative overflow edge", math.MinInt64 / 2, 3, 1, 0},

		// Negative ration overflow cases
		{"negative ration positive amount overflow", math.MaxInt64, -2, 1, 0},
		{"negative ration negative amount overflow", math.MinInt64, -2, 1, 0},

		// Safe cases that should NOT overflow
		{"safe large positive", math.MaxInt64 / 2, 1, 1, math.MaxInt64 / 2},
		{"safe large negative", math.MinInt64 / 2, 1, 1, math.MinInt64 / 2},
		{"safe with division", math.MaxInt64, 1, 2, math.MaxInt64 / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.Allocate(tt.amount, tt.ration, tt.scale)

			if got != tt.want {
				t.Errorf("Allocate(%d, %d, %d) = %d; want %d", tt.amount, tt.ration, tt.scale, got, tt.want)
			}
		})
	}
}

func TestAbsolute(t *testing.T) {
	c := NewCalculator()
	tests := []struct {
		name   string
		amount Amount
		want   Amount
	}{
		{"positive", 100, 100},
		{"negative", -100, 100},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Absolute(tt.amount); got != tt.want {
				t.Errorf("Absolute(%d) = %d; want %d", tt.amount, got, tt.want)
			}
		})
	}
}

func TestNegative(t *testing.T) {
	c := NewCalculator()
	tests := []struct {
		name   string
		amount Amount
		want   Amount
	}{
		{"positive", 100, -100},
		{"negative", -100, -100},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Negative(tt.amount); got != tt.want {
				t.Errorf("Negative(%d) = %d; want %d", tt.amount, got, tt.want)
			}
		})
	}
}

func TestRound(t *testing.T) {
	c := NewCalculator()
	tests := []struct {
		name     string
		calc     *Calculator
		amount   Amount
		exponent int
		want     Amount
	}{
		{"round down", c, 1234, 2, 1200},
		{"round up", c, 1251, 2, 1300},
		{"round down (half)", c, 1250, 2, 1200},
		{"negative round down", c, -1234, 2, -1200},
		{"negative round up", c, -1251, 2, -1300},
		{"nil calculator", nil, 100, 2, 0},
		{"zero amount", c, 0, 2, 0},
		{"exponent zero", c, 1234, 0, 1234},
		{"negative exponent -1", c, 1234, -1, 1234},
		{"negative exponent -2", c, 1251, -2, 1251},
		{"exponent 18 boundary", c, 1234567890123456789, 18, 1000000000000000000},
		{"exponent 19 overflow protection", c, 1234567890123456789, 19, 1234567890123456789},
		{"exponent 20 overflow protection", c, 1234567890123456789, 20, 1234567890123456789},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.calc.Round(tt.amount, tt.exponent); got != tt.want {
				t.Errorf("Round(%d, %d) = %d; want %d", tt.amount, tt.exponent, got, tt.want)
			}
		})
	}
}

func TestCalculatorNilReceiverBranches(t *testing.T) {
	var c *Calculator

	if got := c.Absolute(10); got != 0 {
		t.Fatalf("(*Calculator)(nil).Absolute(10) = %d, want 0", got)
	}

	if got := c.Negative(10); got != 0 {
		t.Fatalf("(*Calculator)(nil).Negative(10) = %d, want 0", got)
	}

	got, err := c.SafeMultiply(2, 3)

	if err != nil || got != 0 {
		t.Fatalf("(*Calculator)(nil).SafeMultiply(2,3) = (%d,%v), want (0,nil)", got, err)
	}
}

func TestCalculatorSafeMultiplyMethodDelegates(t *testing.T) {
	c := NewCalculator()

	got, err := c.SafeMultiply(2, 3, 4)

	if err != nil {
		t.Fatalf("SafeMultiply() unexpected error: %v", err)
	}

	if got != 24 {
		t.Fatalf("SafeMultiply() = %d, want 24", got)
	}

	_, err = c.SafeMultiply(math.MinInt64, -1)

	if err == nil {
		t.Fatal("SafeMultiply(MinInt64, -1) expected overflow error")
	}
}
