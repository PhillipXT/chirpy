package main

import (
    "log"
    "net/http"
    "time"

    "github.com/PhillipXT/chirpy/internal/auth"
)

func (cfg *Config) refreshToken(w http.ResponseWriter, r *http.Request) {
    type response struct {
        Token string `json:"token"`
    }

    bearer, err := auth.GetBearerToken(r.Header)
    if err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "No refresh token found", err)
        return
    }

    log.Printf("[tokens.go] Refresh token found: %s", bearer)

    token, err := cfg.db.GetRefreshToken(r.Context(), bearer)
    if token.RevokedAt.Valid || token.ExpiresAt.Before(time.Now()) {
        log.Printf("Revoked (%v) (%v) : Expired (%v)", token.RevokedAt, token.RevokedAt.Valid, token.ExpiresAt)
        writeErrorResponse(w, http.StatusUnauthorized, "Refresh token is not valid", err)
        return
    }

    new_token, err := auth.MakeJWT(token.UserID, cfg.secret, time.Hour)
    if err != nil {
        writeErrorResponse(w, http.StatusUnauthorized, "Couldn't create new token", err)
        return
    }

    writeResponse(w, http.StatusOK, response { Token: new_token, })
}

func (cfg *Config) revokeToken(w http.ResponseWriter, r *http.Request) {
    bearer, err := auth.GetBearerToken(r.Header)
    if err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "No refresh token found", err)
        return
    }

    err = cfg.db.RevokeRefreshToken(r.Context(), bearer)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error revoking refresh token", err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
