package main

import (
	"encoding/json"
	"log"
	"net/http"
) 
	
type parameters struct {
	Body string `json:"body"`
}

type returnValid struct {
	Valid bool `json:"valid"`
}

type returnError struct {
	Error string `json:"error"`
}

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request)  {


    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)

    if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	const maxchirpslen = 140
	if len(params.Body)	> maxchirpslen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		} else {
		respondWithJSON(w, http.StatusOK, returnValid{Valid: true})
	}
}



