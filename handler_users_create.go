package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/Aleksy37/chirpy/internal/auth"
	"github.com/Aleksy37/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password string `json:"-"`
}


func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request)  {
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

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Hashing password", err)
	}
	user, err  := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user", err)
		return
	}
		
	respondWithJSON(w, http.StatusCreated, response{
		User : User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	})
		
}