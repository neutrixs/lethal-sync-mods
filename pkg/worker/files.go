package web_worker

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type nginxContent struct {
    Name string `json:"name"`
    Type string `json:"type"` //either file or directory
    Mtime string `json:"mtime"`
}

func FetchURLs(link string) ([]string, error) {
    files := []string{}
    var data []nginxContent

    res, err := http.Get(link)
    if err != nil {
        return nil, err
    }

    defer res.Body.Close()

    rawdata, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(rawdata, &data)
    if err != nil {
        return nil, err
    }

    for _, content := range data {
        if content.Type == "file" {
            parsedURL, err := url.Parse(link)
            if err != nil {
                return nil, err
            }

            parsedURL.Path = path.Join(parsedURL.Path, content.Name)

            files = append(files, parsedURL.String())
        } else {
            parsedURL, err := url.Parse(link)
            if err != nil {
                return nil, err
            }

            parsedURL.Path = path.Join(parsedURL.Path, content.Name)
            childURL := parsedURL.String()

            children, err := FetchURLs(childURL)
            if err != nil {
                return nil, err
            }

            files = append(files, children...)
        }
    }

    return files, nil
}

func TrimURL(full string, base string) (string, error) {
    fullPath, err := url.Parse(full)
    if err != nil {
        return "", err
    }

    basePath, err := url.Parse(base)
    if err != nil {
        return "", err
    }

    newPath := strings.TrimPrefix(fullPath.Path, basePath.Path)
    newPath = strings.TrimPrefix(newPath, "/")

    return newPath, nil
}