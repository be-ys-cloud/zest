package services

import (
	"github.com/be-ys-cloud/zest/internal/database"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"os"
	"strings"
)

func DeletePackage(packageName string) (err error) {

	// Get all locations matching
	locations, err := database.GetFileLocation(packageName)
	if err != nil {
		return err
	}

	// Remove all indexes file
	last := packageName[strings.LastIndex(packageName, "/")+1:]
	for _, j := range locations {
		err = os.Remove(configuration.DataStorage + j + "_partials/" + last + "_index")
		if err != nil {
			return err
		}
	}

	// Remove file himself
	_ = os.Remove(configuration.DataStorage + packageName)

	// Clear database : remove file and associated_file locations
	_ = database.DeleteFile(packageName)

	return nil
}
