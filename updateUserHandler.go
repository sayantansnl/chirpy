package main

import (
	"encoding/json"
	"net/http"

	"github.com/sayantansnl/chirpy/internal/auth"
	"github.com/sayantansnl/chirpy/internal/database"
)

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "malformed token", err)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to validate token", err)
		return
	}

	params := request{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to hash password", err)
		return
	}

	user, err := cfg.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userId,
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to update user", err)
		return
	}

	updatedUser := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusOK, updatedUser)
}