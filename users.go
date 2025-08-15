package main

import (
    "net/http"
    "encoding/json"
    "log"
    "time"

    "github.com/google/uuid"
    "github.com/PhillipXT/chirpy/internal/auth"
    "github.com/PhillipXT/chirpy/internal/database"
)

type User struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email string `json:"email"`
    Password string `json:"-"`
    Token string `json:"token"`
    RefreshToken string `json:"refresh_token"`
    IsChirpyRed bool `json:"is_chirpy_red"`
}

func (cfg *Config) createUser(w http.ResponseWriter, r *http.Request) {

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
        writeErrorResponse(w, http.StatusInternalServerError, "Error decoding JSON", err)
        return
    }

    hash, err := auth.HashPassword(params.Password)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error hashing password", err)
        return
    }

    db_params := database.CreateUserParams {
        Password: hash,
        Email: params.Email,
    }

    user, err := cfg.db.CreateUser(r.Context(), db_params)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error creating user", err)
        return
    }

    u := User {
        ID: user.ID,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        Email: user.Email,
        IsChirpyRed: user.IsChirpyRed,
    }

    log.Printf("Created user: %s\n", u.Email)

    writeResponse(w, http.StatusCreated, response { User: u })
}

func (cfg *Config) updateUser(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Email string `json:"email"`
        Password string `json:"password"`
    }

    type response struct {
        User
    }

    bearer, err := auth.GetBearerToken(r.Header)
    if err != nil {
        writeErrorResponse(w, http.StatusUnauthorized, "No access token found", err)
        return
    }

    id, err := auth.ValidateJWT(bearer, cfg.secret)
    if err != nil {
        writeErrorResponse(w, http.StatusUnauthorized, "Invalid token", err)
        return
    }

    params := parameters {}

    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(&params)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Couldn't decode JSON", err)
        return
    }

    hash, err := auth.HashPassword(params.Password)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Couldn't hash password", err)
        return
    }

    update_params := database.UpdateUserParams {
        ID: id,
        Email: params.Email,
        Password: hash,
    }

    user, err := cfg.db.UpdateUser(r.Context(), update_params)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Couldn't update user", err)
        return
    }

    u := User {
        ID: user.ID,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        Email: user.Email,
        IsChirpyRed: user.IsChirpyRed,
    }

    writeResponse(w, http.StatusOK, response { User: u, })
}
