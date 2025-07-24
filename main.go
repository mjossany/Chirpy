package main

import (
	"log"
	"net/http"
)

func main() {
	serverMux := http.NewServeMux()

	const port = "8080"
	const filepathRoot = "."

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}

	serverMux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
