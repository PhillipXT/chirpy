package main

import (
    "net/http"
    "encoding/json"
    "log"
    "time"

    "github.com/google/uuid"
)

type User struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email string `json:"email"`
}

func (cfg *Config) createUser(w http.ResponseWriter, r *http.Request) {

    type parameters struct {
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

    user, err := cfg.db.CreateUser(r.Context(), params.Email)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error creating user", err)
        return
    }

    u := User {
        ID: user.ID,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        Email: user.Email,
    }

    log.Printf("Created user: %s\n", u.Email)

    writeResponse(w, http.StatusCreated, response { User: u })
}
