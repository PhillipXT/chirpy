package auth

import (
    "errors"
    "net/http"
    "strings"
)

func GetAPIKey(headers http.Header) (string, error) {

    header := headers.Get("Authorization")
    if header == "" {
        return "", errors.New("authorization token missing")
    }

    split := strings.Split(header, " ")
    if len(split) < 2 || split[0] != "ApiKey" {
        return "", errors.New("malformed authorization header")
    }

    return split[1], nil

}
