package gpg

import (
	"fmt"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"os"
	"os/exec"
	"time"
)

func InlineSign(privateKeyFileName string, readPass string, fileToSign string) (err error) {
	// Thanks to https://security.stackexchange.com/questions/104149/make-signed-file-from-content-file-and-its-detached-signature

	tempFolder := fmt.Sprintf("%s%d/", configuration.TempStorage, time.Now().UnixNano())

	temporaryFile := tempFolder + "TemporaryInfileSignFile"
	gpgSignatureFile := temporaryFile + ".gpg"

	if err = os.MkdirAll(tempFolder, 0777); err != nil {
		return err
	}

	if err = DetachedSign(privateKeyFileName, readPass, fileToSign, gpgSignatureFile); err != nil {
		return err
	}

	if err = exec.Command("bash", "-c", "echo '-----BEGIN PGP SIGNED MESSAGE-----\nHash: SHA512\n' > "+temporaryFile).Run(); err != nil {
		return err
	}

	if err = exec.Command("bash", "-c", "cat "+fileToSign+" >> "+temporaryFile).Run(); err != nil {
		return err
	}

	if err = exec.Command("bash", "-c", "echo '\r' >> "+temporaryFile).Run(); err != nil {
		return err
	}

	if err = exec.Command("bash", "-c", "cat "+gpgSignatureFile+" >> "+temporaryFile).Run(); err != nil {
		return err
	}

	if err = exec.Command("bash", "-c", "mv "+temporaryFile+" "+fileToSign).Run(); err != nil {
		return err
	}

	if err = os.RemoveAll(tempFolder); err != nil {
		return err
	}

	return nil
}
