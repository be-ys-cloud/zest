package database

import (
	"database/sql"
	"github.com/kisielk/sqlstruct"
	"github.com/be-ys-cloud/zest/internal/structures"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
)

func GetAllFiles() ([]structures.File, error) {
	var files []structures.File

	rows, err := configuration.Database.Query("SELECT * FROM file ORDER BY id ASC;")
	if err != nil {
		return []structures.File{}, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var file structures.File
		err := sqlstruct.Scan(&file, rows)
		if err == nil {
			files = append(files, file)
		}
	}

	return files, nil
}
