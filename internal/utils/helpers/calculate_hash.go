package helpers

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"log"
	"os"
	"strconv"
)

// CalculateHash returns file hashes : MD5, SHA1, SHA256, SHA512
func CalculateHash(file string) map[string]string {

	d, e := os.ReadFile(file)
	if e != nil {
		log.Fatal(e)
	}

	md5sum := md5.Sum(d)
	sha1sum := sha1.Sum(d)
	sha256sum := sha256.Sum256(d)
	sha512sum := sha512.Sum512(d)

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fs, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	return map[string]string{
		"MD5Sum": hex.EncodeToString(md5sum[:]),
		"SHA1":   hex.EncodeToString(sha1sum[:]),
		"SHA256": hex.EncodeToString(sha256sum[:]),
		"SHA512": hex.EncodeToString(sha512sum[:]),
		"Size":   strconv.FormatInt(fs.Size(), 10),
	}
}
