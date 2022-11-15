package database

import (
	"database/sql"
	"github.com/kisielk/sqlstruct"
	"github.com/be-ys-cloud/zest/internal/structures"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
)

func GetFileLocation(fileName string) ([]string, error) {
	rows, err := configuration.Database.Query("SELECT file_location.filePackage AS filePackage FROM file_location INNER JOIN file on file.id = file_location.fileId WHERE file.fileName = ?;", fileName)
	if err != nil {
		return []string{}, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var data structures.FileLocation

	var location []string

	for rows.Next() {
		err := sqlstruct.Scan(&data, rows)
		if err != nil {
			return []string{}, err
		}
		location = append(location, data.FilePackage)
	}

	return location, nil
}
