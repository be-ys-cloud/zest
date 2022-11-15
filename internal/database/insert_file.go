package database

import (
	"github.com/sirupsen/logrus"
	"os"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
)

func InsertFile(path string, packages []string) error {
	fileStats, err := os.Stat(configuration.DataStorage + path)
	size := 0
	if err == nil {
		size = int(fileStats.Size())
	}

	row, err := configuration.Database.Exec("INSERT INTO file(fileName, fileSize) VALUES(?, ?);", path, size)

	if err != nil {
		logrus.Warnln(err.Error())
		return err
	}

	id, _ := row.LastInsertId()

	for _, pack := range packages {
		_, _ = configuration.Database.Exec("INSERT INTO file_location(fileId, filePackage) VALUES(?,?);", id, pack)
	}

	return nil
}
