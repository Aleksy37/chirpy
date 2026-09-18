package main

import (
	"net/http"
	"time"
	"github.com/Aleksy37/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request)  {
	

	type response struct {
		Token string `json:"token"`
	}
	
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Header Error", err)
		return
	}
	
	refreshTokenRecord, err := cfg.db.FindRefreshToken(r.Context(), refreshTokenString)
	if err  != nil {
		respondWithError(w, http.StatusUnauthorized, "This token has no record", err)
		return
	}

	currentTime := time.Now()
	if refreshTokenRecord.ExpiresAt.Before(currentTime) || refreshTokenRecord.ExpiresAt.Equal(currentTime) {
		respondWithError(w, http.StatusUnauthorized, "This token has expired", nil)
		return
	}
	if refreshTokenRecord.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "This token has been revoked", nil)
		return
	}

	accessToken, err := auth.MakeJWT(refreshTokenRecord.UserID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return 
	}
		
	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}