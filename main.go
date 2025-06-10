package main

import (
	"fmt"
	"net/http"
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
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirp)

	//admin
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.reset)

	fmt.Println("Server started")
	server.ListenAndServe()
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
