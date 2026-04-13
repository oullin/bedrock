package events_test

// testOrderCreated is a sample event for testing struct-based event dispatch.
type testOrderCreated struct {
	OrderID string
}

// testOrderShipped is a sample event for testing struct-based event dispatch.
type testOrderShipped struct {
	OrderID   string
	TrackingN string
}

// testUserRegistered is a sample event for testing struct-based event dispatch.
type testUserRegistered struct {
	UserID string
	Email  string
}
