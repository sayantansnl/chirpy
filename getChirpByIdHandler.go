package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirpByIdHandler(w http.ResponseWriter, r *http.Request) {
	val := r.PathValue("id")

	id, err := uuid.Parse(val)
	if err != nil {
		respondWithError(w, http.StatusFailedDependency, "unable to parse id of type uuid.UUID", err)
	}

	chirp, err := cfg.queries.GetChirpById(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to find chirp", err)
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}