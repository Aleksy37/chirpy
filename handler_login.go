package main

import (
	"encoding/json"
	"net/http"
	"github.com/Aleksy37/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request)  {
	type parameter struct {
		Password string `json:"password"`
		Email string `json:"email"`

	}

	type response struct {
		User
	}
	
    decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	user, err  := cfg.db.LoginFindUser(r.Context(), params.Email,)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not find user with that email", err)
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusInternalServerError, "Incorrect email or password", err)
		return
	}
		
	respondWithJSON(w, http.StatusOK, response{
		User : User{
			ID: user.ID,
			Email: user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
		
}