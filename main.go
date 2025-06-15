package main

import (
	"database/sql"
	"fmt"
	"github/breyting/http/internal/database"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Errorf("error connecting to the db")
	}
	dbQueries := database.New(db)

	serveMux := http.NewServeMux()

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	apiCfg := apiConfig{
		queries: dbQueries,
	}

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
