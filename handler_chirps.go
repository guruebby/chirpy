package main

import (
        "net/http"
        "encoding/json"
        "time"
	"fmt"

	"github.com/guruebby/chirpy/internal/database"
        "github.com/google/uuid"
)

type Chirp struct {
        ID        uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Body      string    `json:"body"`
	UserID	  uuid.UUID `json:"user_id"`

}

func validateChirp(body string) (string, error) {

        const maxChirpLength = 140
        if len(body) > maxChirpLength {
                return "", fmt.Errorf("Chirp is too long")
        }

        badWords := map[string]struct{}{
                "kerfuffle": {},
                "sharbert":  {},
                "fornax":    {},
        }
        cleaned := getCleanedBody(body, badWords)

        return cleaned, nil
}



func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
        type parameters struct {
                Body   string `json:"body"`
		UserID string `json:"user_id"`
        }
        type response struct {
                Chirp
        }

        decoder := json.NewDecoder(r.Body)
        params := parameters{}
        err := decoder.Decode(&params)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
                return
        }

	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't provide clean Chirp", err)
	}

	userID, err := uuid.Parse(params.UserID)
	if err != nil {
    		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
    		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Couldn't create Chirp", err)
                return
        }
	respondWithJSON(w, http.StatusCreated, response{
                Chirp: Chirp{
                        ID:        chirp.ID,
                        CreatedAt: chirp.CreatedAt,
                        UpdatedAt: chirp.UpdatedAt,
                        Body:      chirp.Body,
			UserID:	   chirp.UserID,
                },
        })
}
