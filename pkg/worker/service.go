package web_worker

import (
	"errors"
	"io"
	"net/http"
)

var ErrorNotOkResponse = errors.New("Not OK HTTP Response")

type Worker struct {
	Total int64
}

func (w *Worker) Download(url string, dest ...io.Writer) error {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}

	length := res.ContentLength
	w.Total = length

	go io.Copy(io.MultiWriter(dest...), res.Body)

	return nil
}