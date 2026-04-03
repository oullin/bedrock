package calculator

import (
	"errors"
	"math"
	"testing"

	"github.com/gollin/packages/support/money/exception"
)

func TestRationEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		amount Amount
		ration int64
		want   int64
	}{
		{name: "zero ration", amount: 10, ration: 0, want: 0},
		{name: "zero amount", amount: 0, ration: 5, want: 0},
		{name: "min int with negative one overflows", amount: math.MinInt64, ration: -1, want: 0},
		{name: "negative one flips sign", amount: 99, ration: -1, want: -99},
		{name: "positive overflow", amount: math.MaxInt64, ration: 2, want: 0},
		{name: "negative overflow", amount: math.MinInt64 / 2, ration: -3, want: 0},
		{name: "safe negative ration", amount: -5, ration: -2, want: 10},
		{name: "safe positive ration", amount: 7, ration: 3, want: 21},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ration(tt.amount, tt.ration); got != tt.want {
				t.Fatalf("Ration(%d, %d) = %d, want %d", tt.amount, tt.ration, got, tt.want)
			}
		})
	}
}

func TestSafeMultiply(t *testing.T) {
	tests := []struct {
		name        string
		initial     int64
		multipliers []int64
		want        int64
		wantErr     bool
	}{
		{
			name:        "simple multiplication",
			initial:     10,
			multipliers: []int64{2, 3},
			want:        60,
			wantErr:     false,
		},
		{
			name:        "multiplication with negatives",
			initial:     10,
			multipliers: []int64{-1, 2, -3},
			want:        60,
			wantErr:     false,
		},
		{
			name:        "multiplication by zero",
			initial:     100,
			multipliers: []int64{5, 0, 10},
			want:        0,
			wantErr:     false,
		},
		{
			name:        "positive overflow",
			initial:     math.MaxInt64,
			multipliers: []int64{2},
			want:        0,
			wantErr:     true,
		},
		{
			name:        "negative overflow",
			initial:     math.MinInt64,
			multipliers: []int64{2},
			want:        0,
			wantErr:     true,
		},
		{
			name:        "MinInt64 times -1 overflow",
			initial:     math.MinInt64,
			multipliers: []int64{-1},
			want:        0,
			wantErr:     true,
		},
		{
			name:        "-1 times MinInt64 overflow",
			initial:     -1,
			multipliers: []int64{math.MinInt64},
			want:        0,
			wantErr:     true,
		},
		{
			name:        "safe multiplication with -1",
			initial:     100,
			multipliers: []int64{-1},
			want:        -100,
			wantErr:     false,
		},
		{
			name:        "sequential safe multiplications",
			initial:     10,
			multipliers: []int64{2, 3, 4},
			want:        240,
			wantErr:     false,
		},
		{
			name:        "large but safe multiplication",
			initial:     1000000,
			multipliers: []int64{1000, 9},
			want:        9000000000,
			wantErr:     false,
		},
		{
			name:        "overflow in middle of sequence",
			initial:     math.MaxInt64 / 2,
			multipliers: []int64{2, 2},
			want:        0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeMultiply(tt.initial, tt.multipliers...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SafeMultiply() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, exception.ErrOverflow) {
				t.Fatalf("SafeMultiply() unexpected error type: %v", err)
			}
			if got != tt.want {
				t.Fatalf("SafeMultiply(%d, %v) = %d, want %d", tt.initial, tt.multipliers, got, tt.want)
			}
		})
	}
}

func TestSafeMultiply_CoveragePaths(t *testing.T) {
	t.Run("zero short-circuit", func(t *testing.T) {
		got, err := SafeMultiply(10, 0, 5)
		if err != nil {
			t.Fatalf("SafeMultiply() unexpected error: %v", err)
		}
		if got != 0 {
			t.Fatalf("SafeMultiply() = %d, want 0", got)
		}
	})

	t.Run("overflow MinInt64 times -1", func(t *testing.T) {
		_, err := SafeMultiply(math.MinInt64, -1)
		if !errors.Is(err, exception.ErrOverflow) {
			t.Fatalf("SafeMultiply() error = %v, want ErrOverflow", err)
		}
	})

	t.Run("overflow -1 times MinInt64", func(t *testing.T) {
		_, err := SafeMultiply(-1, math.MinInt64)
		if !errors.Is(err, exception.ErrOverflow) {
			t.Fatalf("SafeMultiply() error = %v, want ErrOverflow", err)
		}
	})

	t.Run("overflow positive multiplier", func(t *testing.T) {
		_, err := SafeMultiply(math.MaxInt64, 2)
		if !errors.Is(err, exception.ErrOverflow) {
			t.Fatalf("SafeMultiply() error = %v, want ErrOverflow", err)
		}
	})

	t.Run("overflow negative multiplier less than -1", func(t *testing.T) {
		_, err := SafeMultiply(math.MinInt64, -2)
		if !errors.Is(err, exception.ErrOverflow) {
			t.Fatalf("SafeMultiply() error = %v, want ErrOverflow", err)
		}
	})

	t.Run("safe negative multiplier", func(t *testing.T) {
		got, err := SafeMultiply(10, -2)
		if err != nil {
			t.Fatalf("SafeMultiply() unexpected error: %v", err)
		}
		if got != -20 {
			t.Fatalf("SafeMultiply() = %d, want -20", got)
		}
	})

	t.Run("safe -1 multiplier", func(t *testing.T) {
		got, err := SafeMultiply(5, -1)
		if err != nil {
			t.Fatalf("SafeMultiply() unexpected error: %v", err)
		}
		if got != -5 {
			t.Fatalf("SafeMultiply() = %d, want -5", got)
		}
	})
}
