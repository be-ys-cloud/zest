package controllers

import (
	"encoding/base64"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(401)
			return
		}

		authData := strings.Split(r.Header.Get("Authorization"), " ")[len(strings.Split(r.Header.Get("Authorization"), " "))-1]
		authDataDecoded, err := base64.StdEncoding.DecodeString(authData)

		if err != nil {
			w.WriteHeader(403)
			return
		}

		parts := strings.Split(string(authDataDecoded), ":")

		if parts[0] != "admin" || parts[1] != configuration.AdminPassword {
			w.WriteHeader(403)
			return
		}

		next.ServeHTTP(w, r)
	})
}
