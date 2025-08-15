package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/PhillipXT/chirpy/internal/auth"

    "github.com/google/uuid"
)

func (cfg *Config) upgradeUser(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Event string `json:"event"`
        Data struct {
            UserID uuid.UUID `json:"user_id"`
        } `json:"data"`
    }

    api_token, err := auth.GetAPIKey(r.Header)
    if err != nil {
        writeErrorResponse(w, http.StatusUnauthorized, "Invalid token", err)
        return
    }

    if api_token != cfg.polka_key {
        writeErrorResponse(w, http.StatusUnauthorized, "Invalid api key", err)
        return
    }

    params := parameters{}

    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(&params)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error decoding JSON", err)
        return
    }

    log.Printf("Params: %v => Event (%s) ID (%v)", params, params.Event, params.Data.UserID)

    if params.Event != "user.upgraded" {
        w.WriteHeader(http.StatusNoContent)
        return
    }

    err = cfg.db.UpgradeUser(r.Context(), params.Data.UserID)
    if err != nil {
        writeErrorResponse(w, http.StatusNotFound, "Error upgrading user", err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
