package controllers

import (
	"github.com/be-ys-cloud/zest/internal/services"
	"net/http"
)

func Metrics(w http.ResponseWriter, _ *http.Request) {
	data, err := services.GetMetrics()

	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, _ = w.Write([]byte(data))

}
