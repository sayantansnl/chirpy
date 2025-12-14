package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sayantansnl/chirpy/internal/auth"
)

func (cfg *apiConfig) webhooksHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	params := request{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to decode parameters", err)
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to get API key", err)
		return
	}

	if apiKey != cfg.key {
		respondWithError(w, http.StatusUnauthorized, "mismatched keys", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusFailedDependency, "unable to parse into uuid.UUID", err)
		return
	}
	
	if err := cfg.queries.UpdateChirpyRedByID(r.Context(), userID); err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}