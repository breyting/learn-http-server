package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	serveMux := http.NewServeMux()

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	apiCfg := apiConfig{}

	//website
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	//api
	serveMux.HandleFunc("GET /api/healthz", healthz)
	serveMux.HandleFunc("GET /api/metrics", apiCfg.metrics)
	serveMux.HandleFunc("POST /api/reset", apiCfg.reset)

	fmt.Println("Server started")
	server.ListenAndServe()
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func healthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metrics(w http.ResponseWriter, req *http.Request) {
	res := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(res))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
}
