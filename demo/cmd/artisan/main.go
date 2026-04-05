package main

import (
	"os"

	"github.com/gollin/demo/bootstrap"
)

func main() {
	app, err := bootstrap.New()
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	os.Exit(app.HandleCommand(os.Args[1:], os.Stdout, os.Stderr))
}
