package main

import (
	"context"
	"github.com/be-ys-cloud/zest/internal/controllers"
	"github.com/be-ys-cloud/zest/internal/schedulers"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/scheduler"
	"net/http"
	"time"
)

// main starts our Web Server.
func main() {

	// Start new workers
	ctx := context.Background()
	worker := scheduler.NewScheduler()
	worker.Add(ctx, schedulers.UpdatePackagesFiles, time.Hour*8)
	worker.Add(ctx, schedulers.Cleaner, time.Hour*24+time.Minute*40)

	// Update Packages files on boot
	go schedulers.UpdatePackagesFiles(nil)

	// Start router
	m := mux.NewRouter()
	m.Path("/metrics").HandlerFunc(controllers.Metrics).Methods(http.MethodGet)
	m.PathPrefix("/http").HandlerFunc(controllers.GetData).Methods(http.MethodGet)
	m.PathPrefix("/pass").HandlerFunc(controllers.GetData).Methods(http.MethodGet)

	api := m.PathPrefix("/admin").Subrouter()
	api.Use(controllers.AuthMiddleware)
	api.Path("/stop").HandlerFunc(controllers.StopServer).Methods(http.MethodGet)
	api.Path("/cleanup").HandlerFunc(controllers.Cleanup).Methods(http.MethodGet)
	api.PathPrefix("/packages/").HandlerFunc(controllers.DeletePackage).Methods(http.MethodDelete)

	logrus.Fatalln(http.ListenAndServe(":"+configuration.Port, m))
}
