package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	progressbar "github.com/schollz/progressbar/v3"

	"github.com/neutrixs/lethal-sync-mods/internal/api"
	webworker "github.com/neutrixs/lethal-sync-mods/pkg/worker"
)

const URL = "https://fs.neutrixs.my.id/enjoy/modpack.zip"

func main() {
	worker := webworker.Worker{}

	err := worker.Download(URL)
	if err != nil {
		log.Fatalln(err)
	}

	// Real impressive, windows
	defer func() {
		time.Sleep(time.Second)
		os.Remove(worker.File.Name())
	}()
	defer worker.File.Close()

	prevSize := int64(0)
	bar := progressbar.DefaultBytes(worker.Total, "Downloading")

	for true {
		stat, err := worker.File.Stat()
		if err != nil {
			log.Fatalln(err)
		}

		current := stat.Size()
		bar.Add64(current - prevSize)

		prevSize = current

		if current == worker.Total {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	os.RemoveAll(path.Join(wd, "BepInEx"))

	fmt.Println("Extracting...")
	err = api.Unzip(worker.File, wd)
	fmt.Println("Done!")
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}