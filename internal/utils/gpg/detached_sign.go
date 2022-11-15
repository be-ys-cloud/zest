package gpg

import (
	"crypto"
	"github.com/be-ys-cloud/zest/internal/utils/gpg/configless"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"os"
)

func DetachedSign(privateKeyFileName string, readPass string, fileToSign string, signatureFile string) (err error) {

	var signer *openpgp.Entity
	if signer, err = configless.ReadPrivateKeyFile(privateKeyFileName, readPass); err != nil {
		return err
	}

	var message *os.File
	if message, err = os.Open(fileToSign); err != nil {
		return err
	}

	var w *os.File
	if w, err = os.Create(signatureFile); err != nil {
		return err
	}

	config := packet.Config{DefaultHash: crypto.SHA512}

	if err = openpgp.ArmoredDetachSign(w, signer, message, &config); err != nil {
		return err
	}

	if _, err = w.Write([]byte("\n")); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	if err = message.Close(); err != nil {
		return err
	}

	return nil
}
