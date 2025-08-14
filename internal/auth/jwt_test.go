package auth

import (
    "net/http"
    "testing"
    "time"

    "github.com/google/uuid"
)

// To run tests:  go test ./...
func TestJWT(t *testing.T) {

    userID := uuid.New()
    token, _ := MakeJWT(userID, "secret", time.Hour)

    cases := []struct {
        name string
        tokenString string
        tokenSecret string
        targetID uuid.UUID
        shouldError bool
    } {
        {
            name: "Valid token",
            tokenString: token,
            tokenSecret: "secret",
            targetID: userID,
            shouldError: false,
        },
        {
            name: "Invalid token",
            tokenString: "wrong.token",
            tokenSecret: "secret",
            targetID: uuid.Nil,
            shouldError: true,
        },
        {
            name: "Wrong secret",
            tokenString: token,
            tokenSecret: "wrong secret",
            targetID: uuid.Nil,
            shouldError: true,
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            id, err := ValidateJWT(c.tokenString, c.tokenSecret)
            if err != nil && !c.shouldError {
                t.Errorf("Couldn't validate JWT: %s", err)
            } else if id != c.targetID {
                t.Errorf("ValidateJWT() error: %v, want %v", id, c.targetID)
            }
        })
    }
}

func TestGetBearerToken(t *testing.T) {
    cases := []struct {
        name        string
        headers     http.Header
        wantToken   string
        wantError   bool
    } {
        {
            name: "Valid Bearer Token",
            headers: http.Header { "Authorization": []string{"Bearer valid_token"} },
            wantToken: "valid_token",
            wantError: false,
        }, {
            name: "Missing Authorization Header",
            headers: http.Header {},
            wantToken: "",
            wantError: true,
        }, {
            name: "Malformed Authorization Header",
            headers: http.Header { "Authorization": []string{"InvalidBearer token"} },
            wantToken: "",
            wantError: true,
        },
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            gotToken, err := GetBearerToken(c.headers)
            if (err != nil) != c.wantError {
                t.Errorf("GetBearerToken() error: %v (wantError %v)", err, c.wantError)
            } else if gotToken != c.wantToken {
                t.Errorf("GetBearerToken() gotToken: %v, wanted %v", gotToken, c.wantToken)
            }
        })
    }
}
