package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sayantansnl/chirpy/internal/auth"
	"github.com/sayantansnl/chirpy/internal/database"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to decode parameters", err)
		return
	}

	user, err := cfg.queries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", nil)
		return
	}

	t, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusBadGateway, "unable to generate refresh token", err)
		return
	}

	rT, err := cfg.queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: t,
		UserID: user.ID,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to insert into database refresh_tokens", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusExpectationFailed, "unable to generate token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: rT.Token,
	})

}