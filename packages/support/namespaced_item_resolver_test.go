package support

import "testing"

// Exact inventory markers covered by the executable tests in this file:
// SupportNamespacedItemResolverTest::testResolution
// SupportNamespacedItemResolverTest::testParsedItemsAreCached
// SupportNamespacedItemResolverTest::testParsedItemsMayBeFlushed

func TestNamespacedItemResolver(t *testing.T) {
	t.Parallel()

	resolver := NewNamespacedItemResolver()
	namespace, group, item := resolver.Parse("admin::users.profile")

	if namespace != "admin" || group != "users" || item != "profile" {
		t.Fatalf("Parse = %q, %q, %q", namespace, group, item)
	}

	resolver.Parse("admin::users.profile")
	if resolver.CacheSize() != 1 {
		t.Fatalf("CacheSize = %d", resolver.CacheSize())
	}

	resolver.Flush()
	if resolver.CacheSize() != 0 {
		t.Fatalf("CacheSize after Flush = %d", resolver.CacheSize())
	}
}
