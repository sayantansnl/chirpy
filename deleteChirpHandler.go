package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sayantansnl/chirpy/internal/auth"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	val := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(val)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid auth token", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired token", err)
		return
	}

	chirp, err := cfg.queries.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "action not allowed", nil)
		return
	}

	if err := cfg.queries.DeleteChirpById(r.Context(), chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}