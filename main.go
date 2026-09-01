package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	port := "8080"
	filepathRoot := "."

	apiCfg := &apiConfig{}
	
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetricCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerMetricReset)

	
	svr := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(svr.ListenAndServe())
}



func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}


func (cfg *apiConfig) handlerMetricCount(w http.ResponseWriter, r *http.Request)  {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		w.Write(fmt.Appendf(nil, "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load()))
	}

