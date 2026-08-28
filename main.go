package main

import (
	"log"
	"net/http"
	"path/filepath"
)


func main() {
	port := "8080"
	filepathRoot := "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	svr := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %son port: %s\n", port)
	log.Fatal(svr.ListenAndServe())
}
