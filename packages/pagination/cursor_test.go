package pagination_test

import (
	"testing"

	cpagination "github.com/bedrock/packages/contracts/pagination"
)

func TestCursorEncodeDecode(t *testing.T) {
	t.Parallel()

	cursor := cpagination.NewCursor(map[string]string{"id": "10"}, true)
	encoded := cursor.Encode()

	decoded, err := cpagination.DecodeCursor(encoded)

	if err != nil {
		t.Fatalf("unexpected error decoding cursor: %v", err)
	}

	if decoded.Parameter("id") != "10" {
		t.Errorf("expected parameter id=10, got %q", decoded.Parameter("id"))
	}

	if !decoded.PointsToNextItems() {
		t.Error("expected cursor to point to next items")
	}
}

func TestCursorParameters(t *testing.T) {
	t.Parallel()

	cursor := cpagination.NewCursor(map[string]string{"id": "10", "name": "foo"}, true)
	params := cursor.Parameters()

	if params["id"] != "10" {
		t.Errorf("expected id=10, got %q", params["id"])
	}

	if params["name"] != "foo" {
		t.Errorf("expected name=foo, got %q", params["name"])
	}
}

func TestCursorParameter(t *testing.T) {
	t.Parallel()

	cursor := cpagination.NewCursor(map[string]string{"id": "10"}, true)

	if cursor.Parameter("id") != "10" {
		t.Errorf("expected id=10, got %q", cursor.Parameter("id"))
	}
}
