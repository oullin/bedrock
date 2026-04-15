package echo_test

import (
	"testing"

	"github.com/bedrock/packages/broadcastclient"
)

// TestEventFormatterFormat is a direct port of the six test cases in
// tests/util/event-formatter.test.ts from the Upstream BroadcastClient test suite.
func TestEventFormatterFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		namespace string
		event     string
		want      string
	}{
		{
			name:      "prepends namespace and replaces dot separators with backslashes",
			namespace: "App.Events",
			event:     "Users.UserCreated",
			want:      "App\\Events\\Users\\UserCreated",
		},
		{
			name:      "does not prepend namespace when event starts with a dot",
			namespace: "App.Events",
			event:     ".App\\Users\\UserCreated",
			want:      "App\\Users\\UserCreated",
		},
		{
			name:      "does not prepend namespace when event starts with a backslash",
			namespace: "App.Events",
			event:     "\\App\\Users\\UserCreated",
			want:      "App\\Users\\UserCreated",
		},
		{
			name:      "does not replace dot separators when event starts with a dot",
			namespace: "App.Events",
			event:     ".users.created",
			want:      "users.created",
		},
		{
			name:      "does not replace dot separators when event starts with a backslash",
			namespace: "App.Events",
			event:     "\\users.created",
			want:      "users.created",
		},
		{
			name:      "does not prepend namespace when none is set",
			namespace: "",
			event:     "Users.UserCreated",
			want:      "Users\\UserCreated",
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := broadcastclient.NewEventFormatter(tc.namespace)
			got := f.Format(tc.event)

			if got != tc.want {
				t.Fatalf("Format(%q) with namespace %q = %q; want %q",
					tc.event, tc.namespace, got, tc.want)
			}
		})
	}
}
