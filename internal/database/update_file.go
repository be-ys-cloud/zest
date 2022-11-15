package database

import (
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
)

func UpdateFile(path string) error {

	_, err := configuration.Database.Exec("UPDATE file SET lastDownloaded=CURRENT_TIMESTAMP, nbDownload=nbDownload+1 WHERE fileName=?;", path)

	return err
}
