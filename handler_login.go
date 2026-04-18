package main

import (
        "net/http"
        "encoding/json"

	"github.com/guruebby/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

        decoder := json.NewDecoder(r.Body)
        params := parameters{}
        err := decoder.Decode(&params)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
                return
        }

        user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
        if err != nil {
                respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
                return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
    		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
    		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:	   user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:	   user.Email,
		},
	})
}


