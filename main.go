package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/neutrixs/lethal-sync-mods/internal/api"
	"github.com/neutrixs/lethal-sync-mods/internal/util"
	WebWorker "github.com/neutrixs/lethal-sync-mods/pkg/worker"
)

func main() {
	worker := WebWorker.Worker{}
	worker.Download("https://fs.neutrixs.my.id/enjoy/jawir-modpack.zip")

	for true {
		progress, err := worker.GetProgress()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("\r%s/%s", util.FormatByteSize(progress), util.FormatByteSize(worker.Total))
		time.Sleep(100 * time.Millisecond)

		if progress == worker.Total {
			break
		}
	}

	wd, _ := os.Getwd()
	path := path.Join(wd, ".testing")

	err := api.Unzip(worker.File, path)
	if err != nil {
		fmt.Println(err)
	}

	worker.File.Close()
}