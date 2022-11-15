package controllers

import (
	"github.com/be-ys-cloud/zest/internal/services"
	"net/http"
)

// StopServer gracefully stops a server
func StopServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(202)
	services.StopServer()
}
