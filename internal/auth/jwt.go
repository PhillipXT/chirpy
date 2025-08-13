package auth

import (
    "errors"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/golang-jwt/jwt/v5"
)

const (
    token_issuer = "chirpy"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

    signingKey := []byte(tokenSecret)

    claims := jwt.RegisteredClaims {
        ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
        IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
        Issuer: token_issuer,
        Subject: userID.String(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

    claimsStruct := jwt.RegisteredClaims{}

    token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
        return []byte(tokenSecret), nil
    })
    if err != nil {
        return uuid.Nil, err
    }

    userID, err := token.Claims.GetSubject()
    if err != nil {
        return uuid.Nil, err
    }

    issuer, err := token.Claims.GetIssuer()
    if err != nil {
        return uuid.Nil, err
    } else if issuer != token_issuer {
        return uuid.Nil, errors.New("invalid token issuer")
    }

    id, err := uuid.Parse(userID)
    if err != nil {
        return uuid.Nil, fmt.Errorf("invalid user id: %w", err)
    }

    return id, nil
}
