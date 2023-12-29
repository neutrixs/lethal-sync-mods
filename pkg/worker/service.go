package web_worker

import (
	"errors"
	"io"
	"net/http"
	"os"
)

var ErrorNotOkResponse = errors.New("Not OK HTTP Response")

type Worker struct {
	File *os.File
	Total int64
}

func (w *Worker) GetProgress() (int64, error) {
	stat, err := w.File.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

func (w *Worker) Download(url string) error {
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

	file, err := os.CreateTemp("", "lcsync")
	if err != nil {
		return err
	}

	w.File = file

	go io.Copy(w.File, res.Body)

	return nil
}