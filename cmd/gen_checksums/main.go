package main

import (
	"encoding/json"
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

		hashData, err := api.GetChecksum(fullPath, file)
		if err != nil {
			log.Fatalln(err)
		}

		checksums = append(checksums, hashData)
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