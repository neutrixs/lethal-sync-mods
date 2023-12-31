package api

import (
	"os"
	"path"
	"strings"
)



func GetFiles(dir string) (filePaths []string, err error) {
	files := []string{}
	d, err := os.Open(dir)
	if err != nil {
		return files, err
	}

	stat, err := d.Stat()
	if err != nil {
		return files, err
	}

	// i don't think this will ever be true, but just in case
	if !stat.IsDir() {
		return files, nil
	}

	listFiles, err := d.ReadDir(-1)
	if err != nil {
		return files, err
	}

	for _, file := range listFiles {
		if file.IsDir() {
			newPath := path.Join(dir, file.Name())
			children, err := GetFiles(newPath)
			if err != nil {
				return files, err
			}

			files = append(files, children...)
		} else {
			files = append(files, path.Join(dir, file.Name()))
		}
	}

	return files, nil
}

func GetFilesWithoutTheDir(dir string) (filePaths []string, err error) {
	files := []string{}

	data, err := GetFiles(dir)
	if err != nil {
		return files, err
	}

	for _, file := range data {
		newName := strings.TrimPrefix(file, dir)
		newName = strings.TrimPrefix(newName, "/")

		files = append(files, newName)
	}

	return files, nil
}

type Checksum struct {
	Name string `json:"name"`
	Sha256 string `json:"sha256"`
}