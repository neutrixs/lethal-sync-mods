package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/IGLOU-EU/go-wildcard/v2"

	"github.com/neutrixs/lethal-sync-mods/internal/api"
	web_worker "github.com/neutrixs/lethal-sync-mods/pkg/worker"
	"github.com/schollz/progressbar/v3"
)

const baseURL = "https://lc.neutrixs.my.id"
var filesWL = []string{
	"winhttp.dll",
	"doorstop_config.ini",
	"BepInEx/*",
}

func main() {
	fmt.Println("Verifying checksums...")
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	res, err := http.Get(fmt.Sprintf("%s/%s", baseURL, "checksums.txt"))
	if err != nil {
		log.Fatalln(err)
	}

	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var sourceChecksums []api.Checksum

	err = json.Unmarshal(rawData, &sourceChecksums)
	if err != nil {
		log.Fatalln(err)
	}

	var targetChecksums []api.Checksum

	files, err := api.GetFilesRelative(wd)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		match := false
		for _, pattern := range filesWL {
			if wildcard.Match(pattern, file) {
				match = true
				break
			}
		}

		if !match {
			continue
		}

		fullPath := path.Join(wd, file)

		hashData, err := api.GetChecksum(fullPath, file)
		if err != nil {
			log.Fatalln(err)
		}

		targetChecksums = append(targetChecksums, hashData)
	}

	mis, red := api.CompareChecksums(sourceChecksums, targetChecksums)
	fmt.Printf("%d missing, %d redundant files\n", len(mis), len(red))

	if len(red) > 0 {
		fmt.Println("Removing redundant files...")
		for _, filepath := range red {
			absolutePath := path.Join(wd, filepath)
			err = os.Remove(absolutePath)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	for i, filepath := range mis {
		absolutePath := path.Join(wd, filepath)
		dir, filename := path.Split(absolutePath)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		file, err := os.OpenFile(absolutePath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0755)
		if err != nil {
			log.Fatalln(err)
		}

		bar := progressbar.DefaultBytesSilent(1, fmt.Sprintf("(%d/%d) Downloading %s",i+1, len(mis), filename))

		worker := web_worker.Worker{}
		err = worker.Download(fmt.Sprintf("%s/%s", baseURL, filepath), file, bar)

		bar.ChangeMax64(worker.Total)

		for true {
			stat, err := file.Stat()
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("\r%s", bar.String())

			if stat.Size() == worker.Total {
				break
			}

			time.Sleep(10 * time.Millisecond)
		}

		file.Close()
		bar.Finish()
		bar.Exit()
	}
	fmt.Println()
	fmt.Println("Done! Press the Enter Key to exit anytime")
    fmt.Scanln()
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}