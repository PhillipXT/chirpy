package main

import (
    "encoding/json"
    "net/http"

    "github.com/PhillipXT/chirpy/internal/auth"
)

func (cfg *Config) login(w http.ResponseWriter, r *http.Request) {

    type parameters struct {
        Password string `json:"password"`
        Email string `json:"email"`
    }

    type response struct {
        User
    }

    params := parameters{}

    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&params)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error parsing JSON", err)
        return
    }

    user, err := cfg.db.FindUser(r.Context(), params.Email)
    if err != nil {
        writeErrorResponse(w, http.StatusUnauthorized, "Incorrect email or password", err)
        return
    }

    err = auth.CheckPassword(params.Password, user.Password)
    if err != nil {
        writeErrorResponse(w, http.StatusUnauthorized, "Incorrect email or password", err)
        return
    }

    writeResponse(w, http.StatusOK, response {
        User: User {
            ID: user.ID,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
            Email: user.Email,
        },
    })
}
