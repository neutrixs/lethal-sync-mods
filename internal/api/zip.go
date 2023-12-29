package api

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(file *os.File, dest string) error {
	fi, err := file.Stat()
	if err != nil {
		return err
	}

	reader, err := zip.NewReader(file, fi.Size())
	if err != nil {
		return err
	}

	for _, zf := range reader.File {
		path := filepath.Join(dest, zf.Name)

		if zf.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		os.MkdirAll(filepath.Dir(path), os.ModePerm)

		zfr, err := zf.Open()
		if err != nil {
			return err
		}

		defer zfr.Close()

		destFile, err := os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, zf.Mode())
		if err != nil {
			return err
		}

		defer destFile.Close()

		_, err = io.Copy(destFile, zfr)
		if err != nil {
			return err
		}
	}

	return nil
}