package main

import (
	"net/http"

	"github.com/sayantansnl/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed header", err)
		return
	}

	if err := cfg.queries.SetRevokedAt(r.Context(), token); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}