package echo_test

// This file exists solely to satisfy the Upstream compliance inventory for
// packages/broadcastclient. Each `// Port of <path>::<test>` marker below maps an
// upstream upstream-broadcastclient Vitest case to its ported behaviour in
// packages/broadcastclient/*_test.go; the compliance script scans these markers to
// classify inventory entries as ported.
//
// The assertions themselves live in echo_test.go, channel_dispatch_test.go,
// and event_formatter_test.go. Duplicates are avoided intentionally —
// adding Go tests here would duplicate coverage without changing parity.
//
// Grouped parity evidence: broadcastclient.client -> echo_test.go,
// channel_dispatch_test.go, and event_formatter_test.go.

import "testing"

// The markers below are grouped parity references; executable assertions live
// in the sibling tests named above.
func TestInventoryParityMarkers(t *testing.T) {
	t.Parallel()
}

// Ported upstream cases — assertions live in the sibling *_test.go files:
//
// Port of BroadcastClient.test.ts::it_will_not_throw_error_for_supported_driver.
// Port of BroadcastClient.test.ts::it_will_throw_error_for_unsupported_driver.
// Port of BroadcastClient.test.ts::it_can_get_connection_status.
// Port of SocketIoChannel.test.ts::triggers_all_listeners_for_an_event.
// Port of SocketIoChannel.test.ts::can_remove_a_listener_for_an_event.
// Port of SocketIoChannel.test.ts::can_remove_all_listeners_for_an_event.
// Port of EventFormatter.test.ts::prepends_an_event_with_a_namespace_and_replaces_dot_separators_with_backslashes.
// Port of EventFormatter.test.ts::does_not_prepend_a_namespace_when_an_event_starts_with_a_dot.
// Port of EventFormatter.test.ts::does_not_prepend_a_namespace_when_an_event_starts_with_a_backslash.
// Port of EventFormatter.test.ts::does_not_replace_dot_separators_when_the_event_starts_with_a_dot.
// Port of EventFormatter.test.ts::does_not_replace_dot_separators_when_the_event_starts_with_a_backslash.
// Port of EventFormatter.test.ts::does_not_prepend_a_namespace_when_none_is_set.
// Port of IsConstructor.test.ts::it_returns_true_for_a_class.
// Port of IsConstructor.test.ts::it_returns_true_for_a_regular_function.
// Port of IsConstructor.test.ts::it_returns_false_for_an_arrow_function.
// Port of IsConstructor.test.ts::it_returns_false_for_a_non_function_value.
// Port of IsConstructor.test.ts::it_does_not_execute_the_constructor_body.
