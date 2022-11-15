package services

import (
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"os"
	"time"
)

func StopServer() {
	go stopServer()
}

func stopServer() {

	for configuration.Routines.RunningCount() != 0 || configuration.UnpackRoutines.RunningCount() != 0 || configuration.UpdateInProgress {
		time.Sleep(2 * time.Second)
	}

	os.Exit(0)
}
