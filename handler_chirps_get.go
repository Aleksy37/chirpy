package main

import (
	"net/http"
	"sort"

	"github.com/Aleksy37/chirpy/internal/database"
	"github.com/google/uuid"
) 
	

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}
	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return authorID, nil
}

func (cfg *apiConfig) handlerFetchChirps(w http.ResponseWriter, r *http.Request)  {
	authorID, err := authorIDFromRequest(r)
	sortOrder := r.URL.Query().Get("sort")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}

	var chirps []database.Chirp

	if authorID != uuid.Nil {
		chirps, err = cfg.db.FetchChirpsByUser(r.Context(), authorID)
	} else {
		chirps, err = cfg.db.FetchChirps(r.Context())
	}
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
	if sortOrder == "desc" {
		sort.Slice(feed, func(i, j int) bool {return feed[i].CreatedAt.After(feed[j].CreatedAt)})
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