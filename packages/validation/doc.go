// Package validation provides core functionality for validation.
//
// Ref: @bedrock/code-0187
// Ref: @bedrock/code-0388
// accepts map[string]any data, evaluates 80+ built-in rules expressed as
// pipe-delimited strings ("required|email|max:255"), and collects failures
// into a MessageBag.
//
// Quick start:
//
//	v := validation.NewFactory().Make(
//	    map[string]any{"email": "user@example.com", "age": 17},
//	    map[string]any{"email": "required|email", "age": "required|integer|min:18"},
//	    nil, nil,
//	)
//	if v.Fails() {
//	    fmt.Println(v.Errors().All())
//	}
package validation
