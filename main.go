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

	serverMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	serverMux.HandleFunc("/healthz", handleHealthCheck)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
