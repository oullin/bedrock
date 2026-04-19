package protocol_test

import (
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
)

func TestHeaderConstants(t *testing.T) {
	t.Parallel()

	// Verify header names match the Inertia.js protocol spec.
	if protocol.HeaderInertia != "X-Inertia" {
		t.Errorf("HeaderInertia = %q, want %q", protocol.HeaderInertia, "X-Inertia")
	}

	if protocol.HeaderVersion != "X-Inertia-Version" {
		t.Errorf("HeaderVersion = %q, want %q", protocol.HeaderVersion, "X-Inertia-Version")
	}

	if protocol.HeaderPartialComponent != "X-Inertia-Partial-Component" {
		t.Errorf("HeaderPartialComponent = %q, want %q", protocol.HeaderPartialComponent, "X-Inertia-Partial-Component")
	}

	if protocol.HeaderPartialData != "X-Inertia-Partial-Data" {
		t.Errorf("HeaderPartialData = %q, want %q", protocol.HeaderPartialData, "X-Inertia-Partial-Data")
	}

	if protocol.HeaderPartialExcept != "X-Inertia-Partial-Except" {
		t.Errorf("HeaderPartialExcept = %q, want %q", protocol.HeaderPartialExcept, "X-Inertia-Partial-Except")
	}

	if protocol.HeaderLocation != "X-Inertia-Location" {
		t.Errorf("HeaderLocation = %q, want %q", protocol.HeaderLocation, "X-Inertia-Location")
	}

	if protocol.HeaderInfiniteScroll != "X-Inertia-Infinite-Scroll-Merge-Intent" {
		t.Errorf("HeaderInfiniteScroll = %q, want %q", protocol.HeaderInfiniteScroll, "X-Inertia-Infinite-Scroll-Merge-Intent")
	}

	if protocol.HeaderExceptOnceProps != "X-Inertia-Except-Once-Props" {
		t.Errorf("HeaderExceptOnceProps = %q, want %q", protocol.HeaderExceptOnceProps, "X-Inertia-Except-Once-Props")
	}

	if protocol.HeaderReset != "X-Inertia-Reset" {
		t.Errorf("HeaderReset = %q, want %q", protocol.HeaderReset, "X-Inertia-Reset")
	}

	if protocol.HeaderErrorBag != "X-Inertia-Error-Bag" {
		t.Errorf("HeaderErrorBag = %q, want %q", protocol.HeaderErrorBag, "X-Inertia-Error-Bag")
	}

	if protocol.HeaderRedirect != "X-Inertia-Redirect" {
		t.Errorf("HeaderRedirect = %q, want %q", protocol.HeaderRedirect, "X-Inertia-Redirect")
	}
}
