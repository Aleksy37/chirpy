package main

import (
	"net/http"
	"github.com/Aleksy37/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request)  {
	
	
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Header Error", err)
		return
	}
	
	err = cfg.db.RevokeRefreshToken(r.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Token", err)
		return
	}
	
	respondWithJSON(w, http.StatusNoContent, nil)
}