package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	web_worker "github.com/neutrixs/lethal-sync-mods/pkg/worker"
	"github.com/schollz/progressbar/v3"
)

const baseURL = "https://lc.neutrixs.my.id"

func main() {
	fmt.Println("Cleaning up BepInEx folder...")
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	os.RemoveAll(path.Join(wd, "BepInEx"))

	fmt.Println("Fetching files...")

	fileURLs, err := web_worker.FetchURLs(baseURL)
	if err != nil {
		log.Fatalln(err)
	}

	for i, URL := range fileURLs {
		relativePath, _ := web_worker.TrimURL(URL, baseURL)
		absolutePath := path.Join(wd, relativePath)
		dir, filename := path.Split(absolutePath)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		file, err := os.OpenFile(absolutePath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0755)
		if err != nil {
			log.Fatalln(err)
		}

		bar := progressbar.DefaultBytesSilent(1, fmt.Sprintf("(%d/%d) Downloading %s",i+1, len(fileURLs), filename))

		worker := web_worker.Worker{}
		err = worker.Download(URL, file, bar)

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
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}