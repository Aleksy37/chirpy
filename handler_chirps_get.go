package main

import (
	"net/http"

	"github.com/google/uuid"
) 
	



func (cfg *apiConfig) handlerFetchChirps(w http.ResponseWriter, r *http.Request)  {
	chirps, err := cfg.db.FetchChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching chirps", err)
		return
	}
	 feed :=  make([]Chirp, len(chirps))
	 for i, chirp := range chirps {
		feed[i] = Chirp{
			ID:        chirp.ID,
    		CreatedAt: chirp.CreatedAt,
    		UpdatedAt: chirp.UpdatedAt,
    		Body:      chirp.Body,
    		UserID:    chirp.UserID,
		}
	 }

	respondWithJSON(w, http.StatusOK, feed)
}

func (cfg *apiConfig) handlerFetchChirpByID(w http.ResponseWriter, r *http.Request)  {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldnt parse chirpID", err)
	}
	chirp, err := cfg.db.FetchChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp by that ID doesn't exist", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}