package database

import (
	"database/sql"
	"github.com/kisielk/sqlstruct"
	"github.com/be-ys-cloud/zest/internal/structures"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
)

func GetFileByName(fileName string) (structures.File, error) {
	rows, err := configuration.Database.Query("SELECT * FROM file WHERE fileName=? LIMIT 1;", fileName)
	if err != nil {
		return structures.File{}, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var data structures.File

	for rows.Next() {
		err := sqlstruct.Scan(&data, rows)
		if err != nil {
			return structures.File{}, err
		}
	}

	return data, nil
}
