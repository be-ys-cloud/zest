package controllers

import (
	"github.com/be-ys-cloud/zest/internal/services"
	"net/http"
)

// Cleanup inspects all packages in server, and remove ones who are too old
func Cleanup(w http.ResponseWriter, r *http.Request) {
	if err := services.Cleanup(); err != nil {
		w.WriteHeader(500)
	}
}
