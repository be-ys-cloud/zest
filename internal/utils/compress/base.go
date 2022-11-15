package compress

import (
	"errors"
	"strings"
)

// Inflate Create an archive from a plaintext. Format is defined depending on extension. Supported methods are : xz, gz, bz2
func Inflate(sourceFile string, destFile string) error {
	format := strings.Split(destFile, ".")

	switch format[len(format)-1] {
	case "xz":
		return xzInflate(sourceFile, destFile)
	case "gz":
		return gzInflate(sourceFile, destFile)
	case "bz2":
		return bz2Inflate(sourceFile, destFile)
	default:
		return errors.New("extension not supported")
	}
}

// Deflate Decompress an archive and store content as plaintext. Format is defined depending on extension. Supported methods are : xz, gz, bz2
func Deflate(sourceFile string, destFile string) error {
	format := strings.Split(sourceFile, ".")

	switch format[len(format)-1] {
	case "xz":
		return xzDeflate(sourceFile, destFile)
	case "gz":
		return gzDeflate(sourceFile, destFile)
	case "bz2":
		return bz2Deflate(sourceFile, destFile)
	default:
		return errors.New("extension not supported")
	}
}
