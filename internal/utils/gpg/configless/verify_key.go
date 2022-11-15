package configless

import (
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"os"
)

func ReadPrivateKeyFile(filename string, passPrompt string) (e *openpgp.Entity, err error) {
	var krpriv *os.File

	if krpriv, err = os.Open(filename); err != nil {
		return nil, err
	}
	defer krpriv.Close()

	var entityList openpgp.EntityList

	keyFileReader := openpgp.ReadKeyRing
	if _, err = armor.Decode(krpriv); err == nil {
		keyFileReader = openpgp.ReadArmoredKeyRing
	}

	if _, err = krpriv.Seek(0, 0); err != nil {
		return nil, err
	}
	if entityList, err = keyFileReader(krpriv); err != nil {
		return nil, fmt.Errorf("reading %s: %s", filename, err)
	}
	if len(entityList) != 1 {
		return nil, fmt.Errorf("%s must contain only one key", filename)
	}
	e = entityList[0]
	if e.PrivateKey == nil {
		return nil, fmt.Errorf("%s does not contain a private key", filename)
	}
	if e.PrivateKey.Encrypted {
		err = e.PrivateKey.Decrypt([]byte(passPrompt))
	}
	return e, err
}
