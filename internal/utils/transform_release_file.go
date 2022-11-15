package utils

import (
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"github.com/be-ys-cloud/zest/internal/utils/gpg"
	"github.com/sirupsen/logrus"
	"os"
)

func TransformReleaseFile(folderName string, fileFullPath string, body []byte) (err error) {
	//Copy file
	inReleaseFile := configuration.DataStorage + folderName + "/InRelease"

	err = os.MkdirAll(configuration.DataStorage+folderName, 0777)
	if err != nil {
		logrus.Warnln("Could not create directory " + configuration.DataStorage + folderName)
		return
	}

	err = os.WriteFile(configuration.DataStorage+fileFullPath, body, 0777)
	if err != nil {
		logrus.Warnln("Could not write file " + configuration.DataStorage + fileFullPath)
		return
	}

	err = os.WriteFile(inReleaseFile, body, 0777)
	if err != nil {
		logrus.Warnln("Could not write file " + inReleaseFile)
		return
	}

	//Re-create Release.gpg and InRelease file
	err = gpg.DetachedSign(configuration.KeyFile, configuration.KeyPass, configuration.DataStorage+fileFullPath, configuration.DataStorage+fileFullPath+".gpg")
	if err != nil {
		logrus.Warnln("Could not create detached file for " + configuration.DataStorage + fileFullPath)
		return
	}

	err = gpg.InlineSign(configuration.KeyFile, configuration.KeyPass, inReleaseFile)
	if err != nil {
		logrus.Warnln("Could not create inline sign for " + inReleaseFile)
		return
	}

	return
}
