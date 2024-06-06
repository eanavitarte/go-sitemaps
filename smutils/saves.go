package smutils

import (
	"compress/gzip"
	"os"
)

func WriteToFile(filePath string, xmlContent []byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(xmlContent)
	if err != nil {
		return err
	}

	return nil
}

func WriteToGzip(filePath string, xmlContent []byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	zip, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer zip.Close()

	_, err = zip.Write(xmlContent)
	if err != nil {
		return err
	}

	return nil
}
