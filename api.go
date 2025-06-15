package main

import (
	"encoding/json"
	"github/breyting/http/internal/database"
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
}

func healthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateChirp(w http.ResponseWriter, req *http.Request) {
	//Je pourrais DRY ce code facilement
	type parameters struct {
		Body string `json:"body"`
	}

	type error struct {
		Error string `json:"error"`
	}

	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respBody := error{
			Error: "Decode went wrong",
		}

		dat, _ := json.Marshal(respBody)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	if len(params.Body) <= 140 {
		badWords := []string{"kerfuffle", "sharbert", "fornax"}
		text := params.Body

		splitText := strings.Split(text, " ")

		for i, word := range splitText {
			lowerWord := strings.ToLower(word)
			if slices.Contains(badWords, lowerWord) {
				splitText[i] = "****"
			}
		}

		respText := strings.Join(splitText, " ")

		respBody := returnVal{
			CleanedBody: respText,
		}
		dat, _ := json.Marshal(respBody)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
	} else {
		respBody := error{
			Error: "Chirp is too long",
		}

		dat, _ := json.Marshal(respBody)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
	}

}
