package container_test

import (
	"testing"

	"github.com/bedrock/packages/container"
)

func TestTagAndTaggedResolvesAll(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("report-a", func(_ *container.Container) (any, error) {
		return "a", nil
	}, false)

	c.Bind("report-b", func(_ *container.Container) (any, error) {
		return "b", nil
	}, false)

	c.Tag([]string{"report-a", "report-b"}, "reports")

	results := c.Tagged("reports")

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0] != "a" || results[1] != "b" {
		t.Fatalf("unexpected results: %v", results)
	}
}

func TestTaggedReturnsEmptyForUnknownTag(t *testing.T) {
	t.Parallel()

	c := newContainer()

	results := c.Tagged("nonexistent")

	if len(results) != 0 {
		t.Fatalf("expected empty results, got %d", len(results))
	}
}

func TestMultipleTagsOnSameAbstract(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("service", func(_ *container.Container) (any, error) {
		return "service", nil
	}, false)

	c.Tag([]string{"service"}, "tag-a", "tag-b")

	a := c.Tagged("tag-a")
	b := c.Tagged("tag-b")

	if len(a) != 1 || a[0] != "service" {
		t.Fatalf("tag-a: expected [service], got %v", a)
	}

	if len(b) != 1 || b[0] != "service" {
		t.Fatalf("tag-b: expected [service], got %v", b)
	}
}

func TestTagAccumulatesAbstracts(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("a", func(_ *container.Container) (any, error) { return "a", nil }, false)
	c.Bind("b", func(_ *container.Container) (any, error) { return "b", nil }, false)

	c.Tag([]string{"a"}, "services")
	c.Tag([]string{"b"}, "services")

	results := c.Tagged("services")

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestTaggedSkipsUnresolvableAbstracts(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("good", func(_ *container.Container) (any, error) {
		return "good", nil
	}, false)

	c.Tag([]string{"good", "missing"}, "mixed")

	results := c.Tagged("mixed")

	if len(results) != 1 || results[0] != "good" {
		t.Fatalf("expected [good], got %v", results)
	}
}
