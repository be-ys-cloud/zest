package controllers

import (
	"github.com/be-ys-cloud/zest/internal/services"
	"net/http"
	"strings"
)

// GetData is the only controller used to poll data from debian caches.
func GetData(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimSuffix(r.URL.Path, "/")
	url = strings.TrimPrefix(url, "/")

	data, err := services.GetData(url)

	if err != nil {
		if err.Error() == "html" {
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(500)
		return
	}

	_, _ = w.Write(data)

}
