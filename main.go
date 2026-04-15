package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const port = "8080"
	const filepathRoot = "."
	var apiCfg apiConfig
	mux := http.NewServeMux()

	rmPrefixHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(rmPrefixHandler))
	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleMetricReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Start serving from %s on port %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerHealthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	log.Println("Server status being checked")
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (c *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf(`
	<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
	</html>`, c.fileserverHits.Load())
	w.Write([]byte(hits))
	log.Println("Metrics checked: ", c.fileserverHits.Load())
}

func (c *apiConfig) handleMetricReset(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	c.fileserverHits.Store(0)
	w.Write([]byte("Server metrics reset\n"))
	log.Println("Metrics reset")
}

func handlerValidate(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	msg := chirp{}
	err := decoder.Decode(&msg)
	if err != nil {
		respondWithError(w, 500, "Error decoding chirp", err)
		return
	}

	if len(msg.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	success := struct {
		Valid bool `json:"valid"`
	}{Valid: true}
	respondWithJSON(w, 200, success)
	log.Printf("Chirp length good: %d", len(msg.Body))
}

func respondWithError(w http.ResponseWriter, code int, msg string, errMSG error) {
	log.Printf("Encountered: %s\n		Sending: %d-%s", errMSG, code, msg)

	errResp := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	respondWithJSON(w, code, errResp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Marshaling JSON: %s", err)
		respondWithError(w, 500, "Internal server error", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(data)

}
