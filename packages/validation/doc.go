// Package validation is a 1:1 Go port of upstream/framework 13.x
// src/Framework/Validation.  It provides a rule-based input validator that
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
