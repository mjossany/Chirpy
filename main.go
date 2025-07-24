package main

import "net/http"

func main() {
	serverMux := http.NewServeMux()

	const port = "8080"

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}

	serverMux.Handle("/", http.FileServer(http.Dir(".")))

	srv.ListenAndServe()
}
