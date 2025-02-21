package auth

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func GetAuth(username, password string) *http.BasicAuth {
	if password == "" {
		return nil
	}

	if username == "" {
		username = "token"
	}

	return &http.BasicAuth{
		Username: username,
		Password: password,
	}
}
