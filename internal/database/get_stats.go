package database

import (
	"database/sql"
	"github.com/kisielk/sqlstruct"
	"github.com/sirupsen/logrus"
	"github.com/be-ys-cloud/zest/internal/structures"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
)

func GetStats() (structures.Statistics, error) {
	rows, err := configuration.Database.Query("SELECT SUM(fileSize * (nbDownload -1)) as savedBytes, COUNT(*) as nbFiles, SUM(fileSize) as totalSizeOnDisk FROM file;")
	if err != nil {
		logrus.Warnln(err.Error())
		return structures.Statistics{}, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var data structures.Statistics

	for rows.Next() {
		err := sqlstruct.Scan(&data, rows)
		if err != nil {
			logrus.Warnln(err.Error())
			return structures.Statistics{}, err
		}
	}

	return data, nil
}
