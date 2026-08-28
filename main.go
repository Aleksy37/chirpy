package main

import (
	"log"
	"net/http"
)



func main() {
	port := "8080"
	filepathRoot := "."
	
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", handlerReadiness)
	
	svr := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(svr.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request)  {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(200)
			w.Write([]byte(http.StatusText(http.StatusOK)))
		}