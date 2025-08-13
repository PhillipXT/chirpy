package auth

import (
    "golang.org/x/crypto/bcrypt"
)

const (
    cost = 10
)

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func CheckPassword(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
