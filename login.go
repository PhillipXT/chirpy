package main

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/PhillipXT/chirpy/internal/auth"
    "github.com/PhillipXT/chirpy/internal/database"
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

    expirationTime := time.Hour

    token, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error creating token", err)
        return
    }

    refresh_token := auth.MakeRefreshToken()

    refresh_params := database.CreateRefreshTokenParams {
        Token: refresh_token,
        UserID: user.ID,
        ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
    }

    _, err = cfg.db.CreateRefreshToken(r.Context(), refresh_params)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error creating refresh token", err)
        return
    }

    writeResponse(w, http.StatusOK, response {
        User: User {
            ID: user.ID,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
            Email: user.Email,
            Token: token,
            RefreshToken: refresh_token,
        },
    })
}
