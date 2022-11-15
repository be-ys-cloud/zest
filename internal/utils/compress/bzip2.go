package compress

import (
	"github.com/dsnet/compress/bzip2"
	"io"
	"os"
)

func bz2Inflate(sourceFile string, destFile string) error {
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}

	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}

	destWriter, err := bzip2.NewWriter(dest, &bzip2.WriterConfig{})
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

func bz2Deflate(sourceFile string, destFile string) error {
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}

	r, err := bzip2.NewReader(source, &bzip2.ReaderConfig{})
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
