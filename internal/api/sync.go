package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/IGLOU-EU/go-wildcard/v2"
	"github.com/neutrixs/lethal-sync-mods/constants"
	web_worker "github.com/neutrixs/lethal-sync-mods/pkg/worker"
	"github.com/schollz/progressbar/v3"
)

// target must be absolute path
func SyncToClient(source string, target string, whitelist []string, ignorelist []string) error {
	fmt.Println("Verifying checksums...")
	
	csURLData, err := url.Parse(source)
	if err != nil { return err }

	csURLData.Path = path.Join(csURLData.Path, "checksums.txt")
	csURL := csURLData.String()

	res, err := http.Get(csURL)
	if err != nil { return err }
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	rawChecksums, err := io.ReadAll(res.Body)

	var sourceChecksums []Checksum
	var targetChecksums []Checksum

	err = json.Unmarshal(rawChecksums, &sourceChecksums)
	if err != nil { return err }

	files, err := GetFilesRelative(target)
	if err != nil { return err }

	for _, file := range files {
		var match bool
		for _, pattern := range constants.ModsWhitelist {
			if wildcard.Match(pattern, file) {
				match = true
				break
			}
		}

		for _, pattern := range constants.ModsIgnore {
			if wildcard.Match(pattern, file) {
				match = false
				break
			}
		}

		if !match { continue }

		hash, err := GetChecksum(path.Join(target, file), file)
		if err != nil { return err }

		targetChecksums = append(targetChecksums, hash)
	}

	missing, redundant := CompareChecksums(sourceChecksums, targetChecksums)
	fmt.Printf("%d missing, %d redundant files\n", len(missing), len(redundant))

	if len(redundant) > 0 {
		fmt.Println("Removing redundant files...")
		for _, filepath := range redundant {
			err = os.Remove(path.Join(target, filepath))
			if err != nil { return err }
		}
	}

	for i, filepath := range missing {
		absolutePath := path.Join(target, filepath)
		absoluteDir, filename := path.Split(absolutePath)
		
		sourceURLData, err := url.Parse(source)
		if err != nil { return err }
		sourceURLData.Path = path.Join(sourceURLData.Path, filepath)

		err = os.MkdirAll(absoluteDir, 0755)
		if err != nil { return err }

		file, err := os.OpenFile(absolutePath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0755)
		if err != nil { return err }

		formatting := fmt.Sprintf("(%d/%d) Downloading %s", i+1, len(missing), filename)
		bar := progressbar.DefaultBytesSilent(0, formatting)
		
		worker := web_worker.Worker{}
		err = worker.Download(sourceURLData.String(), file, bar)
		if err != nil {
			file.Close()
			bar.Close()
			return err
		}
		
		bar.ChangeMax64(worker.Total)

		for true {
			stat, err := file.Stat()
			if err != nil {
				file.Close()
				bar.Close()
				return err
			}

			fmt.Printf("\r%s", bar.String())
			// >= --> better be safe than sorry
			if stat.Size() >= worker.Total {
				file.Close()
				bar.Close()
				break
			}

			time.Sleep(10 * time.Millisecond)
		}

		file.Close()
		bar.Close()
	}
	fmt.Println()
	return nil
}