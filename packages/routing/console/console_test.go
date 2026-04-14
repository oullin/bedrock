package console

import (
	"strings"
	"testing"
)

func TestControllerMakeCommand(t *testing.T) {
	out := ControllerMakeCommand{}.Render("controllers", "UserController")

	if !strings.Contains(out, "package controllers") {
		t.Errorf("missing package: %q", out)
	}

	if !strings.Contains(out, "type UserController struct") {
		t.Errorf("missing type: %q", out)
	}
}

func TestMiddlewareMakeCommand(t *testing.T) {
	out := MiddlewareMakeCommand{}.Render("middleware", "Authenticate")

	if !strings.Contains(out, "package middleware") {
		t.Errorf("missing package: %q", out)
	}

	if !strings.Contains(out, "type Authenticate struct") {
		t.Errorf("missing type: %q", out)
	}
}
