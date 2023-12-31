package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/neutrixs/lethal-sync-mods/internal/api"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	files, err := api.GetFilesRelative(wd)
	if err != nil {
		log.Fatalln(err)
	}

	checksums := []api.Checksum{}

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

		hashData := api.Checksum{
			Name: file,
			Sha256: fmt.Sprintf("%x",string(h.Sum(nil))),
		}
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