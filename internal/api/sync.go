package api

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"time"

	web_worker "github.com/neutrixs/lethal-sync-mods/pkg/worker"
	"github.com/schollz/progressbar/v3"
)

// target must be absolute path
func SyncToClient(source string, target string, whitelist []string, ignorelist []string) error {
    fmt.Println("Verifying checksums...")

    var sourceChecksums []Checksum
    var targetChecksums []Checksum

    sourceChecksums, err := GetRemoteChecksums(source, "checksums.txt")
    if err != nil {
        log.Println(err)
        return err
    }

    targetChecksums, err = GetChecksums(target, whitelist, ignorelist)
    if err != nil {
        log.Println(err)
        return err
    }

    missing, redundant := CompareChecksums(sourceChecksums, targetChecksums)
    fmt.Printf("%d missing, %d redundant files\n", len(missing), len(redundant))

    if len(redundant) > 0 {
        fmt.Println("Removing redundant files...")
        for _, filepath := range redundant {
            err = os.Remove(path.Join(target, filepath))
            if err != nil {
                log.Println(err)
                return err
            }
        }
    }

    for i, filepath := range missing {
        absolutePath := path.Join(target, filepath)
        absoluteDir, filename := path.Split(absolutePath)
        
        sourceURLData, err := url.Parse(source)
        if err != nil {
            log.Println(err)
            return err
        }
        sourceURLData.Path = path.Join(sourceURLData.Path, filepath)

        err = os.MkdirAll(absoluteDir, 0755)
        if err != nil {
            log.Println(err)
            return err
        }

        file, err := os.OpenFile(absolutePath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0755)
        if err != nil {
            log.Println(err)
            return err
        }

        formatting := fmt.Sprintf("(%d/%d) Downloading %s", i+1, len(missing), filename)
        bar := progressbar.DefaultBytesSilent(0, formatting)
        
        worker := web_worker.Worker{}
        err = worker.Download(sourceURLData.String(), file, bar)
        if err != nil {
            file.Close()
            bar.Close()
            log.Println(err)
            return err
        }
        
        bar.ChangeMax64(worker.Total)

        for {
            stat, err := file.Stat()
            if err != nil {
                file.Close()
                bar.Close()
                log.Println(err)
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

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}