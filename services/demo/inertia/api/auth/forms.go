package auth

import (
	"net/http"
	"strings"

	"github.com/bedrock/packages/inertia/protocol"
)

type loginForm struct {
	Email    string
	Password string
	Remember bool
}

func newLoginForm(r *http.Request) loginForm {
	remember := r.FormValue("remember")

	return loginForm{
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: r.FormValue("password"),
		Remember: remember == "on" || remember == "true" || remember == "1",
	}
}

func (f loginForm) validate() protocol.ValidationErrors {
	errors := protocol.ValidationErrors{}

	if strings.TrimSpace(f.Email) == "" {
		errors["email"] = "Enter the demo email address."
	}

	if strings.TrimSpace(f.Password) == "" {
		errors["password"] = "Enter the demo password."
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
