package main

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/PhillipXT/chirpy/internal/auth"
)

func (cfg *Config) login(w http.ResponseWriter, r *http.Request) {

    type parameters struct {
        Password string `json:"password"`
        Email string `json:"email"`
        Expiry int `json:"-"`
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

    expirationTime := time.Hour
    if params.Expiry > 0 && params.Expiry < 3600 {
        expirationTime = time.Duration(params.Expiry) * time.Second
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

    token, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error creating token", err)
        return
    }

    writeResponse(w, http.StatusOK, response {
        User: User {
            ID: user.ID,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
            Email: user.Email,
            Token: token,
        },
    })
}
