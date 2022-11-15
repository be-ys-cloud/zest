package services

import (
	"github.com/be-ys-cloud/zest/internal/database"
	"github.com/be-ys-cloud/zest/internal/structures"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func Cleanup() error {
	logrus.Infoln("Cache cleanup requested at", time.Now().String(), "is starting.")

	files, err := database.GetAllFiles()
	if err != nil {
		return err
	}

	packagelist := map[string][]structures.File{}

	for _, d := range files {
		//https://www.debian.org/doc/manuals/debian-faq/pkg-basics.en.html
		fileName := d.FileName[strings.LastIndex(d.FileName, "/")+1:]
		folderName := d.FileName[0:strings.LastIndex(d.FileName, "/")]

		folderName += strings.Split(fileName, "_")[0] + "_" + strings.TrimSuffix(strings.Split(fileName, "_")[2], ".deb")
		d.Version = strings.Split(fileName, "_")[1]

		packagelist[folderName] = append(packagelist[folderName], d)
	}

	logrus.Infoln(len(packagelist), "packages found in database.")

	days := 60
	if d, err := strconv.Atoi(configuration.MaxRetentionTime); err == nil {
		days = d
	}

	for _, u := range packagelist {
		for _, v := range u[0 : len(u)-1] {
			if v.LastDownloaded.AddDate(0, 0, days).Unix() < time.Now().Unix() {
				logrus.Infoln("Package", v.FileName, "will be removed, as it was not downloaded for", days, "days, and a newer version exists in cache.")

				err = DeletePackage(v.FileName)
				if err != nil {
					logrus.Warnln("Error while removing package", v.FileName, "! Error was :", err.Error())
				} else {
					logrus.Infoln("Package", v.FileName, "successfully removed.")
				}
			}
		}
	}

	logrus.Infoln("Cleanup done.")

	return nil
}
