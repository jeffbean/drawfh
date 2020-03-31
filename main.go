package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jeffbean/drawfh/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %v", port)
	}
	addr := net.JoinHostPort("", port)

	mux := http.NewServeMux()
	mux.Handle("/", server.Handler())

	return http.ListenAndServe(addr, mux)
}
