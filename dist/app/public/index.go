package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("DEMO_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Fatal(http.ListenAndServe(addr, http.DefaultServeMux))
}
