package action

// PayInvoice handles paying an outstanding invoice.
//
// In the Go implementation this delegates to the payment provider SDK,
// which is injected at the application level. The action itself is a
// thin orchestration wrapper.
type PayInvoice struct{}
