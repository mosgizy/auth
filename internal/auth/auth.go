package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetToken(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("No authentication info found")
	}

	vals := strings.Split(val," ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "Bearer" {
		return "", errors.New("malformed first part of header: 'Bearer'")
	}

	return vals[1], nil
}