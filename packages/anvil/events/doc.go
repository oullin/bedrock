// Package events provides a Upstream-inspired event dispatcher.
//
// The dispatcher manages listeners for named events and supports wildcard
// patterns, halting dispatch (Until), and subscriber registration. It is
// safe for concurrent use.
//
//	d := events.New()
//	d.Listen([]string{"user.created"}, func(ctx context.Context, event string, payload any) (any, error) {
//	    fmt.Println("user created:", payload)
//	    return nil, nil
//	})
//	d.Dispatch(context.Background(), "user.created", user)
package events
