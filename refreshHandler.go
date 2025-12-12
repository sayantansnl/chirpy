package main

import (
	"net/http"
	"time"

	"github.com/sayantansnl/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed header", err)
		return
	}

	rT, err := cfg.queries.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to get from table refresh_tokens", err)
		return
	}

	if time.Now().UTC().After(rT.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "refresh token expired", nil)
		return
	}

	if rT.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token revoked", nil)
		return
	}

	accessToken, err := auth.MakeJWT(rT.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create access JWT", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}