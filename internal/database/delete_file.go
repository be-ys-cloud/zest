package database

import (
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
)

func DeleteFile(fileName string) error {
	data, err := GetFileByName(fileName)

	if err != nil {
		return err
	}

	_, _ = configuration.Database.Exec("DELETE FROM file_location WHERE fileId=?;", data.Id)
	_, _ = configuration.Database.Exec("DELETE FROM file WHERE id=?;", data.Id)

	return nil
}
