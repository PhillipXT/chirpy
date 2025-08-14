package main

import (
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/PhillipXT/chirpy/internal/auth"
    "github.com/PhillipXT/chirpy/internal/database"
)

const maxChirpLength = 140

type Chirp struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body string `json:"body"`
    UserID uuid.UUID `json:"user_id"`
}

func (cfg *Config) getChirp(w http.ResponseWriter, r *http.Request) {

    type response struct {
        Chirp
    }

    uuid, err := uuid.Parse(r.PathValue("chirpID"))
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error parsing UUID", err)
        return
    }

    data, err := cfg.db.GetChirp(r.Context(), uuid)
    if err != nil {
        writeErrorResponse(w, http.StatusNotFound, "Error retrieving chirp", err)
        return
    }

    log.Printf("Found chirp: %s\n", data.Body)

    chirp := Chirp {
        ID: data.ID,
        CreatedAt: data.CreatedAt,
        UpdatedAt: data.UpdatedAt,
        Body: data.Body,
        UserID: data.UserID,
    }

    writeResponse(w, http.StatusOK, response { Chirp: chirp })
}

func (cfg *Config) getChirps(w http.ResponseWriter, r *http.Request) {

    data, err := cfg.db.GetChirps(r.Context())
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error retrieving chirps", err)
        return
    }

    log.Printf("Retrieved %d chirp(s)\n", len(data))

    chirps := []Chirp{}

    for _, row := range data {
        chirps = append(chirps, Chirp {
            ID: row.ID,
            CreatedAt: row.CreatedAt,
            UpdatedAt: row.UpdatedAt,
            Body: row.Body,
            UserID: row.UserID,
        })
    }

    writeResponse(w, http.StatusOK, chirps)
}

func (cfg *Config) createChirp(w http.ResponseWriter, r *http.Request) {

    type parameters struct {
        Body string `json:"body"`
    }

    type response struct {
        Chirp
    }

    log.Printf("[chirps.go] Headers: %v", r.Header)

    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        writeErrorResponse(w, http.StatusUnauthorized, "Token missing", err)
        return
    }

    log.Printf("[chirps.go] Received token: %v", token)

    id, err := auth.ValidateJWT(token, cfg.secret)
    if err != nil {
        writeErrorResponse(w, http.StatusUnauthorized, "Token invalid", err)
        return
    }

    params := parameters {}

    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(&params)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error decoding JSON", err)
        return
    }

    cleaned, err := validateChirp(&params.Body)
    if err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "Chirp is too long", err)
        return
    }

    chirp_params := database.CreateChirpParams {
        Body: cleaned,
        UserID: id,
    }

    log.Printf("Chirp params: %v", chirp_params)

    fc, err := cfg.db.CreateChirp(r.Context(), chirp_params)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error creating chirp", err)
        return
    }

    o := Chirp {
        ID: fc.ID,
        CreatedAt: fc.CreatedAt,
        UpdatedAt: fc.UpdatedAt,
        Body: fc.Body,
        UserID: fc.UserID,
    }

    log.Printf("[chirps.go] Created chirp: %s\n", o.ID)

    writeResponse(w, http.StatusCreated, response { Chirp: o })
}

func validateChirp(chirp *string) (string, error) {

    log.Printf("[chirps.go] Chirp length: %d\n", len(*chirp))
    if len(*chirp) > maxChirpLength {
        return "", errors.New("Chirp is too long")
    }

    cleaned := filterChirp(*chirp)
    log.Printf("[chirps.go] Filtered chirp: %s\n", cleaned)

    return cleaned, nil
}

func filterChirp(message string) string {
    badWords := map[string]struct{} {
        "kerfuffle": {},
        "sharbert": {},
        "fornax": {},
    }

    words := strings.Split(message, " ")

    for i, word := range words {
        if _, ok := badWords[strings.ToLower(word)]; ok {
            words[i] = "****"
        }
    }

    return strings.Join(words, " ")
}
