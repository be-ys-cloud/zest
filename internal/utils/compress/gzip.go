package compress

import (
	"compress/gzip"
	"io"
	"os"
)

func gzInflate(sourceFile string, destFile string) error {
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}

	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}

	destWriter := gzip.NewWriter(dest)

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

func gzDeflate(sourceFile string, destFile string) error {
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}

	r, err := gzip.NewReader(source)
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
