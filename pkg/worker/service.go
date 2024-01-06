package web_worker

import (
	"errors"
	"io"
	"log"
	"net/http"
)

var ErrorNotOkResponse = errors.New("not OK HTTP Response")

type Worker struct {
    Total int64
}

func (w *Worker) Download(url string, dest ...io.Writer) error {
    request, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        log.Println(err)
        return err
    }

    client := http.Client{}
    res, err := client.Do(request)
    if err != nil {
        log.Println(err)
        return err
    }

    length := res.ContentLength
    w.Total = length

    go io.Copy(io.MultiWriter(dest...), res.Body)

    return nil
}

func init() {
    log.SetFlags(log.LstdFlags | log.Llongfile)
}