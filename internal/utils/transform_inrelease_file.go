package utils

import (
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"github.com/be-ys-cloud/zest/internal/utils/gpg"
	cp "github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

func TransformInReleaseFile(folderName string, fileFullPath string, body []byte) (bodyReturn []byte, err error) {
	//Copy file
	releaseFile := configuration.DataStorage + folderName + "/Release"
	err = os.MkdirAll(configuration.DataStorage+folderName, 0777)
	if err != nil {
		logrus.Warnln("Could not create directory " + configuration.DataStorage + folderName)
		return
	}

	err = os.WriteFile(releaseFile, body, 0777)
	if err != nil {
		logrus.Warnln("Could not create file " + releaseFile)
		return
	}

	//Remove armor around file content
	err = exec.Command("sed", "-i", "-e", "1,3d", releaseFile).Run()
	if err != nil {
		logrus.Warnln("Could not create sed in file " + releaseFile)
		return
	}

	err = exec.Command("bash", "-c", "truncate -s $(grep -b '^-----BEGIN PGP SIGNATURE-----' "+releaseFile+" | grep -o -E '^[0-9]*') "+releaseFile).Run()
	if err != nil {
		logrus.Warnln("Could not grep in file " + releaseFile)
		return
	}

	//Change Signed-By field
	err = ReplaceValidUntilAndSignedBy(releaseFile)
	if err != nil {
		logrus.Warnln("Could not change Signed-By field. Program will continue, but some clients may reject your file.")
	}

	//Re-create Release.gpg and InRelease file
	err = gpg.DetachedSign(configuration.KeyFile, configuration.KeyPass, releaseFile, releaseFile+".gpg")
	if err != nil {
		logrus.Warnln("Could not create detached sign for file " + releaseFile)
		return
	}
	err = cp.Copy(releaseFile, configuration.DataStorage+fileFullPath)
	if err != nil {
		logrus.Warnln("Could not copy file " + releaseFile + " to " + configuration.DataStorage + fileFullPath)
		return
	}

	err = gpg.InlineSign(configuration.KeyFile, configuration.KeyPass, configuration.DataStorage+fileFullPath)
	if err != nil {
		logrus.Warnln("Could not inline sign " + configuration.DataStorage + fileFullPath)
		return
	}

	//Return InRelease content to body variable
	f, err := os.Open(configuration.DataStorage + fileFullPath)
	if err != nil {
		logrus.Warnln("Could not open file " + configuration.DataStorage + fileFullPath)
		return
	}
	bodyReturn, err = io.ReadAll(f)
	if err != nil {
		logrus.Warnln("Could not read file content" + configuration.DataStorage + fileFullPath)
		return
	}

	return
}
