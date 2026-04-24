package spark

import "testing"

func TestStringFromAny(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want string
	}{
		{name: "string", in: "hello", want: "hello"},
		{name: "int", in: 1, want: ""},
		{name: "nil", in: nil, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringFromAny(tt.in); got != tt.want {
				t.Fatalf("stringFromAny(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestIntFromAny(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want int
	}{
		{name: "int", in: 7, want: 7},
		{name: "int64", in: int64(8), want: 8},
		{name: "float64", in: float64(9.9), want: 9},
		{name: "string", in: "10", want: 0},
		{name: "nil", in: nil, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intFromAny(tt.in); got != tt.want {
				t.Fatalf("intFromAny(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestBoolFromAny(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want bool
	}{
		{name: "true", in: true, want: true},
		{name: "false", in: false, want: false},
		{name: "string", in: "true", want: false},
		{name: "nil", in: nil, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := boolFromAny(tt.in); got != tt.want {
				t.Fatalf("boolFromAny(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestBillableConfigsFromMap(t *testing.T) {
	in := map[string]any{
		"typed":     BillableConfig{Model: "typed", TrialDays: 3, DefaultInterval: "yearly"},
		"untyped":   map[string]any{"model": "u", "trial_days": 5, "default_interval": "monthly"},
		"untyped64": map[string]any{"trial_days": int64(10)},
		"ignored":   "not-a-map",
	}

	got := billableConfigsFromMap(in)

	if got["typed"].DefaultInterval != "yearly" {
		t.Fatalf("typed entry DefaultInterval = %q", got["typed"].DefaultInterval)
	}

	if got["untyped"].TrialDays != 5 {
		t.Fatalf("untyped trial days = %d", got["untyped"].TrialDays)
	}

	if got["untyped"].Model != "u" {
		t.Fatalf("untyped model = %q", got["untyped"].Model)
	}

	if got["untyped64"].TrialDays != 10 {
		t.Fatalf("untyped64 trial days = %d", got["untyped64"].TrialDays)
	}

	if _, ok := got["ignored"]; ok {
		t.Fatal("unknown shape should be dropped")
	}
}

func TestFeatureConfigsFromMap(t *testing.T) {
	in := map[string]any{
		"typed":   FeatureConfig{Enabled: true, Options: map[string]any{"k": "v"}},
		"untyped": map[string]any{"enabled": true, "options": map[string]any{"x": 1}},
		"noopts":  map[string]any{"enabled": true},
		"skipped": []string{"not", "a", "map"},
	}

	got := featureConfigsFromMap(in)

	if !got["typed"].Enabled || got["typed"].Options["k"] != "v" {
		t.Fatalf("typed = %#v", got["typed"])
	}

	if !got["untyped"].Enabled || got["untyped"].Options["x"] != 1 {
		t.Fatalf("untyped = %#v", got["untyped"])
	}

	if !got["noopts"].Enabled || got["noopts"].Options == nil {
		t.Fatalf("noopts options should default to empty map: %#v", got["noopts"])
	}

	if _, ok := got["skipped"]; ok {
		t.Fatal("unknown shape should be skipped")
	}
}
