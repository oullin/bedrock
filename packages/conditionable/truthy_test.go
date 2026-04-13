package conditionable_test

import (
	"testing"

	"github.com/bedrock/packages/conditionable"
)

func TestTruthy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"nil is falsy", nil, false},
		{"false is falsy", false, false},
		{"true is truthy", true, true},
		{"zero int is falsy", 0, false},
		{"positive int is truthy", 1, true},
		{"negative int is truthy", -1, true},
		{"zero int8 is falsy", int8(0), false},
		{"nonzero int8 is truthy", int8(1), true},
		{"zero int16 is falsy", int16(0), false},
		{"nonzero int16 is truthy", int16(1), true},
		{"zero int32 is falsy", int32(0), false},
		{"nonzero int32 is truthy", int32(1), true},
		{"zero int64 is falsy", int64(0), false},
		{"nonzero int64 is truthy", int64(1), true},
		{"zero uint is falsy", uint(0), false},
		{"nonzero uint is truthy", uint(1), true},
		{"zero uint8 is falsy", uint8(0), false},
		{"nonzero uint8 is truthy", uint8(1), true},
		{"zero uint16 is falsy", uint16(0), false},
		{"nonzero uint16 is truthy", uint16(1), true},
		{"zero uint32 is falsy", uint32(0), false},
		{"nonzero uint32 is truthy", uint32(1), true},
		{"zero uint64 is falsy", uint64(0), false},
		{"nonzero uint64 is truthy", uint64(1), true},
		{"zero float32 is falsy", float32(0), false},
		{"nonzero float32 is truthy", float32(3.14), true},
		{"zero float64 is falsy", float64(0), false},
		{"nonzero float64 is truthy", float64(3.14), true},
		{"empty string is falsy", "", false},
		{"nonempty string is truthy", "hello", true},
		{"struct is truthy", struct{}{}, true},
		{"nil slice is falsy", ([]int)(nil), false},
		{"empty slice is truthy", []int{}, true},
		{"nil map is falsy", (map[string]int)(nil), false},
		{"empty map is truthy", map[string]int{}, true},
		{"nil pointer is falsy", (*int)(nil), false},
		{"non-nil pointer is truthy", new(int), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := conditionable.Truthy(tt.value); got != tt.want {
				t.Errorf("Truthy(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
