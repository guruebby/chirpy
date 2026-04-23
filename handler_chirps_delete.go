package main

import (
        "net/http"

        "github.com/google/uuid"
	"github.com/guruebby/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
        if err != nil {
                respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
                return
        }

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
        if err != nil {
                respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
                return
        }

	chirpIDString := r.PathValue("chirpID")
        chirpID, err := uuid.Parse(chirpIDString)
        if err != nil {
                respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
                return
        }

        dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
        if err != nil {
                respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
                return
        }

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Chirp not by user", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
