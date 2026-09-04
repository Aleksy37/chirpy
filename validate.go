package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
) 
	
type parameters struct {
	Body string `json:"body"`
}

type returnValid struct {
	CleanedBody string `json:"cleaned_body"`
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
			cleaned := getCleanBody(params.Body, badWords)
			respondWithJSON(w, http.StatusOK, returnValid{CleanedBody: cleaned})
	}
}


var badWords = map[string]struct{}{
"kerfuffle": {},
"sharbert":  {},
"fornax":    {},
}

func getCleanBody(text string, badWords map[string]struct{}) string {
	words := strings.Split(text, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}



