package broadcasting_test

import (
	"testing"

	"github.com/bedrock/packages/broadcasting"
)

func TestPusherChannelNameNormalization(t *testing.T) {
	t.Parallel()

	// UsePusherChannelsNamesTest::testChannelNameNormalization
	// UsePusherChannelsNamesTest::testChannelNameNormalizationSpecialCase
	cases := map[string]string{
		"private-test":                  "test",
		"private-encrypted-test":        "test",
		"presence-test":                 "test",
		"public-test":                   "public-test",
		"private-private-test":          "private-test",
		"private-presence-test":         "presence-test",
		"presence-private-test":         "private-test",
		"presence-presence-test":        "presence-test",
		"private-encrypted-private-123": "private-123",
	}

	for input, want := range cases {
		if got := broadcasting.NormalizeChannelName(input); got != want {
			t.Fatalf("NormalizeChannelName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestPusherChannelNamePatternMatching(t *testing.T) {
	t.Parallel()

	// UsePusherChannelsNamesTest::testChannelNamePatternMatching
	// UsePusherChannelsNamesTest::testChannelNameMatchesPattern
	b := broadcasting.NewBaseBroadcaster()

	if b.ChannelNameMatchesPattern("TestChannel", "Test.{id}") {
		t.Fatal("expected TestChannel not to match Test.{id}")
	}

	if !b.ChannelNameMatchesPattern("orders.42", "orders.{id}") {
		t.Fatal("expected orders.42 to match orders.{id}")
	}
}

func TestPusherIsGuardedChannel(t *testing.T) {
	t.Parallel()

	// UsePusherChannelsNamesTest::testIsGuardedChannel
	cases := map[string]bool{
		"private-test":           true,
		"private-encrypted-test": true,
		"presence-test":          true,
		"test":                   false,
		"public-test":            false,
	}

	for channel, want := range cases {
		if got := broadcasting.IsGuardedChannel(channel); got != want {
			t.Fatalf("IsGuardedChannel(%q) = %v, want %v", channel, got, want)
		}
	}
}
