package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")

	if authorID == "" {
		chirps, err := cfg.queries.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "unable to get chirps", err)
		}

		chirpsArray := []Chirp{}

		for _, chirp := range chirps {
			chirpsArray = append(chirpsArray, Chirp{
				ID: chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body: chirp.Body,
				UserID: chirp.UserID,
			})
		}

		respondWithJSON(w, http.StatusOK, chirpsArray)
	} else {
		userID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "unable to parse into uuid.UUID", err)
			return
		} 
		chirps, err := cfg.queries.GetAllChirpsByUserID(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "unable to get chirps", err)
		}

		chirpsArray := []Chirp{}

		for _, chirp := range chirps {
			chirpsArray = append(chirpsArray, Chirp{
				ID: chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body: chirp.Body,
				UserID: chirp.UserID,
			})
		}

		respondWithJSON(w, http.StatusOK, chirpsArray)
	}
}