package main

import (
	"encoding/json"
	"net/http"
	"github.com/Aleksy37/chirpy/internal/auth"
	"github.com/Aleksy37/chirpy/internal/database"
)


func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request)  {
	type parameter struct {
		Email string `json:"email"`
		Password string `json:"password"`

	}

	type response struct {
		User
	}
	
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
    decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Hashing password", err)
	}
	err = cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: hashedPass,
		ID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
		return
	}

	user, err := cfg.db.LoginFindUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating the user", err)
		return
	}
		
	respondWithJSON(w, http.StatusOK, response{
		User : User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	})
		
}