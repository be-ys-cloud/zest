package compress

import (
	"github.com/ulikunitz/xz"
	"io"
	"os"
)

func xzInflate(sourceFile string, destFile string) error {
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}

	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}

	destWriter, err := xz.NewWriter(dest)
	if err != nil {
		return err
	}

	_, err = io.Copy(destWriter, source)
	if err != nil {
		return err
	}

	err = destWriter.Close()
	if err != nil {
		return err
	}

	err = dest.Close()
	if err != nil {
		return err
	}

	err = source.Close()
	if err != nil {
		return err
	}

	return nil
}

func xzDeflate(sourceFile string, destFile string) error {
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}

	r, err := xz.NewReader(source)
	if err != nil {
		return err
	}

	dataContent, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	err = os.WriteFile(destFile, dataContent, 777)
	if err != nil {
		return err
	}

	err = source.Close()
	if err != nil {
		return err
	}

	return nil
}
