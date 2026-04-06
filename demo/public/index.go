package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bedrock/demo/bootstrap"
)

func main() {
	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err)
	}

	addr := os.Getenv("DEMO_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Fatal(http.ListenAndServe(addr, app))
}
