package reverb_test

import (
	"testing"

	"github.com/bedrock/packages/reverb"
)

func TestValidateOrigin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		origin         string
		allowedOrigins []string
		want           bool
	}{
		{
			name:           "exact match",
			origin:         "example.com",
			allowedOrigins: []string{"example.com"},
			want:           true,
		},
		{
			name:           "wildcard matches subdomain",
			origin:         "sub.example.com",
			allowedOrigins: []string{"*.example.com"},
			want:           true,
		},
		{
			name:           "wildcard does not match apex",
			origin:         "example.com",
			allowedOrigins: []string{"*.example.com"},
			want:           false,
		},
		{
			name:           "star pattern matches anything",
			origin:         "anything.io",
			allowedOrigins: []string{"*"},
			want:           true,
		},
		{
			name:           "empty allowedOrigins allows all",
			origin:         "example.com",
			allowedOrigins: []string{},
			want:           true,
		},
		{
			name:           "multiple patterns first match wins",
			origin:         "foo.example.com",
			allowedOrigins: []string{"*.example.com", "other.com"},
			want:           true,
		},
		{
			name:           "no match returns false",
			origin:         "other.com",
			allowedOrigins: []string{"example.com"},
			want:           false,
		},
		{
			name:           "origin with scheme matches bare hostname pattern",
			origin:         "https://example.com",
			allowedOrigins: []string{"example.com"},
			want:           true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := reverb.ValidateOrigin(tc.origin, tc.allowedOrigins)

			if got != tc.want {
				t.Errorf("ValidateOrigin(%q, %v) = %v, want %v", tc.origin, tc.allowedOrigins, got, tc.want)
			}
		})
	}
}
