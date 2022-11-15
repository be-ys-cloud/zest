package configuration

import (
	"database/sql"
	"fmt"
	"github.com/be-ys-cloud/zest/internal/utils/gpg/configless"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/zenthangplus/goccm"
	"os"
	"strconv"
	"strings"
)

var (
	SupportedExtensions = []string{"xz", "bz2", "gz"}
	KeyFile             = "key.asc"
	KeyPass             = ""
	DataStorage         = "data/"
	TempStorage         = "tmp/"
	Port                = "80"
	KeyName             = ""
	MaxRoutines         = "10"
	AdminPassword       = "12345"
	Routines            = goccm.New(10)
	UnpackRoutines      = goccm.New(10)
	UpdateInProgress    = false
	DatabaseFile        = "database.sql"
	MaxRetentionTime    = "60"
	FreeSpaceThreshold  = "10"
	Database            *sql.DB
)

// init gets information from environment and store them in variables ready to be used in other packages.
func init() {
	mapping := map[string]*string{
		"PORT":                 &Port,
		"TEMP_STORAGE":         &TempStorage,
		"DATA_STORAGE":         &DataStorage,
		"KEY_PASSWORD":         &KeyPass,
		"KEY_FILE":             &KeyFile,
		"MAX_ROUTINES":         &MaxRoutines,
		"ADMIN_PASSWORD":       &AdminPassword,
		"DATABASE_FILE":        &DatabaseFile,
		"MAX_RETENTION_TIME":   &MaxRetentionTime,
		"FREE_SPACE_THRESHOLD": &FreeSpaceThreshold,
	}

	// Get configuration from environment and populate values.
	for i, k := range mapping {
		if value, ok := os.LookupEnv("ZEST_" + i); ok && value != "" {
			*k = value
		}

		if strings.HasSuffix(i, "_STORAGE") && !strings.HasSuffix(*k, "/") {
			*k = fmt.Sprintf("%s/", *k)
		}
	}

	// Read Private Key file and get Identity Name.
	entity, err := configless.ReadPrivateKeyFile(KeyFile, KeyPass)

	if err != nil {
		logrus.Fatalln("Cannot read key file !. Ensure the file exists, and the password you provided is valid.")
	}

	if AdminPassword == "12345" {
		logrus.Warnln("You didn't set an admin password, so we will use 12345. We highly recommend you to change this password by setting the ZEST_ADMIN_PASSWORD environment variable.")
	}

	maxRoutinesInt, err := strconv.Atoi(MaxRoutines)
	if err == nil {
		Routines = goccm.New(maxRoutinesInt)
		UnpackRoutines = goccm.New(maxRoutinesInt)
	} else {
		MaxRoutines = "10"
	}

	KeyName = strings.ToUpper(fmt.Sprintf("%x", entity.PrivateKey.Fingerprint))

	// Attempting to create/read database file, and create tables...
	Database = getDatabaseConnection()
	_, _ = Database.Exec("CREATE TABLE file(`id` INTEGER PRIMARY KEY AUTOINCREMENT, `fileName` VARCHAR(1024) NOT NULL, `lastDownloaded` TIMESTAMP DEFAULT CURRENT_TIMESTAMP, `nbDownload` INTEGER DEFAULT 1, `fileSize` INTEGER NOT NULL);")
	_, _ = Database.Exec("CREATE TABLE file_location(`id` INTEGER PRIMARY KEY AUTOINCREMENT, `fileId` INTEGER, `filePackage` VARCHAR(1024) NOT NULL, FOREIGN KEY(`fileId`) REFERENCES `file`(`id`));")

}

func getDatabaseConnection() *sql.DB {

	var db *sql.DB
	var err error

	db, err = sql.Open("sqlite3", DatabaseFile)

	if err != nil {
		logrus.Error("Unable to connect to database !")
		logrus.Fatal(err.Error())
	}

	return db
}
