package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Serving on port: %s", port)
	log.Fatal(server.ListenAndServe())
}