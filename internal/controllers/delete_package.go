package controllers

import (
	"github.com/be-ys-cloud/zest/internal/services"
	"net/http"
	"strings"
)

// DeletePackage is used to remove a package from server.
func DeletePackage(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimPrefix(r.URL.Path, "/admin/packages/")

	if !strings.HasSuffix(url, ".deb") {
		w.WriteHeader(400)
		return
	}

	if err := services.DeletePackage(url); err != nil {
		w.WriteHeader(500)
		return
	}
}
