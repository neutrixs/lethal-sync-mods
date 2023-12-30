package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

func getFiles(dir string) (filePaths []string, err error) {
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
			children, err := getFiles(newPath)
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

func getFilesWithoutTheDir(dir string) (filePaths []string, err error) {
	files := []string{}

	data, err := getFiles(dir)
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

type checksum struct {
	Name string `json:"name"`
	Sha256 string `json:"sha256"`
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	files, err := getFilesWithoutTheDir(wd)
	if err != nil {
		log.Fatalln(err)
	}

	checksums := []checksum{}

	for _, file := range files {
		fullPath := path.Join(wd, file)
		f, err := os.OpenFile(fullPath, os.O_RDONLY, 0755)
		if err != nil {
			log.Fatalln(err)
		}

		h := sha256.New()

		if _, err := io.Copy(h, f); err != nil {
			f.Close()
			log.Fatalln(err)
		}

		hashData := checksum{file, fmt.Sprintf("%x",string(h.Sum(nil)))}
		checksums = append(checksums, hashData)
		f.Close()
	}

	output, err := os.OpenFile(path.Join(wd, "checksums.txt"), os.O_RDONLY | os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := json.Marshal(checksums)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err = output.Write(data); err != nil {
		output.Close()
		log.Fatalln(err)
	}

	output.Close()
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}