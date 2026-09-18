package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"github.com/joho/godotenv"
	"os"
	"database/sql"
	"github.com/Aleksy37/chirpy/internal/database"
)

import _ "github.com/lib/pq"

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
	secret string
}

func main() {
	
	const port = "8080"
	const filepathRoot = "."

	godotenv.Load()
  	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL Must be set")
	}
	plat := os.Getenv("PLATFORM")
	if plat == "" {
		log.Fatal("PLATFORM must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening db: %s", err)
	}
	dbQueries := database.New(dbConn)

	secret := os.Getenv("SECRET")

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: plat,
		secret: secret,
	}
	
	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerFetchChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerFetchChirpByID)
	
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetricCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerMetricReset)

	
	svr := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(svr.ListenAndServe())
}