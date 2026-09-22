package main

import (
	"net/http"
	"github.com/Aleksy37/chirpy/internal/auth"
	"github.com/google/uuid"
) 

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request)  {
	
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or malformed token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

    chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldnt parse chirpID", err)
		return
	}
	dbChirp, err := cfg.db.FetchChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp by that ID doesn't exist", err)
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "User is not chirp author", err)
		return 
	}
	err = cfg.db.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return 
	}
	w.WriteHeader(http.StatusNoContent)
}
