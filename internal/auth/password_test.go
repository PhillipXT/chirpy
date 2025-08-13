package auth

import (
    "testing"
)

// To run tests:  go test ./...
func TestPassword(t *testing.T) {
    cases := []struct {
        name string
        pw1 string
        pw2 string
        shouldError bool
    } {
        {
            name: "Correct password",
            pw1: "hello",
            pw2: "hello",
            shouldError: false,
        }, {
            name: "Incorrect password",
            pw1: "hello",
            pw2: "goodbye",
            shouldError: true,
        },
        {
            name: "Empty password",
            pw1: "",
            pw2: "hello",
            shouldError: true,
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            hash, err := HashPassword(c.pw1)
            if err != nil {
                t.Errorf("Couldn't hash password: %s", err)
            }
            result := CheckPassword(c.pw2, hash)
            if result != nil && !c.shouldError {
                t.Errorf("Passwords do not match: %s vs %s", c.pw1, c.pw2)
            }
        })
    }
}

