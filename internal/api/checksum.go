package api

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/IGLOU-EU/go-wildcard/v2"
)

type Checksum struct {
    Name string `json:"name"`
    Sha256 string `json:"sha256"`
}

func GetFiles(dir string) (filePaths []string, err error) {
    files := []string{}
    d, err := os.Open(dir)
    if err != nil {
        log.Println(err)
        return files, err
    }

    stat, err := d.Stat()
    if err != nil {
        log.Println(err)
        return files, err
    }

    // i don't think this will ever be true, but just in case
    if !stat.IsDir() {
        log.Println(err)
        return files, nil
    }

    listFiles, err := d.ReadDir(-1)
    if err != nil {
        log.Println(err)
        return files, err
    }

    for _, file := range listFiles {
        if file.IsDir() {
            newPath := path.Join(dir, file.Name())
            children, err := GetFiles(newPath)
            if err != nil {
                log.Println(err)
                return files, err
            }

            files = append(files, children...)
        } else {
            files = append(files, path.Join(dir, file.Name()))
        }
    }

    return files, nil
}

// returns the path of the file relative to dir
func GetFilesRelative(dir string) (filePaths []string, err error) {
    files := []string{}

    data, err := GetFiles(dir)
    if err != nil {
        log.Println(err)
        return files, err
    }

    for _, file := range data {
        newName := strings.TrimPrefix(file, dir)
        newName = strings.TrimPrefix(newName, "/")

        files = append(files, newName)
    }

    return files, nil
}

// missing also means different hash, which needs to be redownloaded.
// returns missing, redundant
func CompareChecksums(source []Checksum, target []Checksum) (mis []string, red []string) {
    missing := []string{}
    redundant := []string{}

    for _, checksum := range source {
        index := slices.IndexFunc(target, func(el Checksum) bool{return el.Name == checksum.Name})
        if index == -1 || target[index].Sha256 != checksum.Sha256 {
            missing = append(missing, checksum.Name)
            continue
        }
    }

    for _, checksum := range target {
        if index := slices.IndexFunc(source, func(el Checksum) bool{return el.Name == checksum.Name}); index == -1 {
            redundant = append(redundant, checksum.Name)
        }
    }

    return missing, redundant
}

//name can also be relative path from anywhere, e.g relative/path/from/somewhere.file
func GetChecksum(filepath string, name string) (Checksum, error) {
    var cs Checksum

    f, err := os.OpenFile(filepath, os.O_RDONLY, 0755)
    if err != nil {
        log.Println(err)
        return cs, err
    }
    defer f.Close()

    h := sha256.New()

    if _, err := io.Copy(h, f); err != nil {
        return cs, err
    }

    cs.Name = name
    cs.Sha256 = fmt.Sprintf("%x",string(h.Sum(nil)))

    return cs, nil
}

func GetChecksums(dir string, whitelist []string, ignore []string) ([]Checksum, error) {
    var result []Checksum

    files, err := GetFilesRelative(dir)
    if err != nil { 
        log.Println(err)
        return result, err
    }

    for _, file := range files {
        var match bool
        for _, pattern := range whitelist {
            if wildcard.Match(pattern, file) {
                match = true
                break
            }
        }

        for _, pattern := range ignore{
            if wildcard.Match(pattern, file) {
                match = false
                break
            }
        }

        if !match { continue }

        hash, err := GetChecksum(path.Join(dir, file), file)
        if err != nil {
            log.Println(err)
            return result, err
        }

        result = append(result, hash)
    }

    return result, nil
}

func GetRemoteChecksums(baseURL string, filename string) ([]Checksum, error) {
    var result []Checksum
    csURLData, err := url.Parse(baseURL)
    if err != nil {
        log.Println(err)
        return result, err
    }

    csURLData.Path = path.Join(csURLData.Path, "checksums.txt")
    csURL := csURLData.String()

    res, err := http.Get(csURL)
    if err != nil {
        log.Println(err)
        return result, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return result, errors.New(res.Status)
    }

    rawChecksums, err := io.ReadAll(res.Body)
    if err != nil {
        log.Println(err)
        return result, err
    }

    err = json.Unmarshal(rawChecksums, &result)
    if err != nil {
        log.Println(err)
        return result, err
    }

    return result, nil
}

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}