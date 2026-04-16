package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/jman2476/httpServersLearn/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening database connection: %s", err)
	}

	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
		platform:       platform,
	}

	mux := http.NewServeMux()
	pathPrefixStrip := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	fsHandler := apiCfg.middlewareMetricsInc(pathPrefixStrip)

	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerNewChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerNewUser)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Start serving from %s on port %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
