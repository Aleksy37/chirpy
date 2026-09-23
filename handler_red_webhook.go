package main

import (
	"encoding/json"
	"net/http"

	"github.com/Aleksy37/chirpy/internal/auth"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerRedWebhook(w http.ResponseWriter, r *http.Request)  {

	type parameter struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
		
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if apiKey != cfg.polka || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Api key", err)
		return
	}
	
    decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userUUID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse provided userID into uuid", err)
		return
	}

	err  = cfg.db.UserChirpyRed(r.Context(), userUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find user with that id", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)	
}