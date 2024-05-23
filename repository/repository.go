package repository

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
	"tt/entity"
)

const timeFormat = time.RFC3339

type Repo struct {
	filePath string
}

func NewRepo(filePath string) *Repo {
	return &Repo{filePath: filePath}
}

func (r *Repo) OpenWindow(t time.Time) error {
	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	record := fmt.Sprintf("%s,", t.Format(timeFormat))
	if _, err := file.WriteString(record); err != nil {
		return err
	}

	return nil
}

func (r *Repo) CloseWindow(t time.Time) error {
	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	record := fmt.Sprintf("%s\n", t.Format(timeFormat))
	if _, err := file.WriteString(record); err != nil {
		return err
	}

	return nil
}

func (r *Repo) Read(from, to time.Time) ([]entity.Window, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var windows []entity.Window
	reader := csv.NewReader(file)
	reader.ReuseRecord = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		w, err := deserialize(record)
		if err != nil {
			return nil, err
		}

		if !w.StartsAt.Before(from) && to.After(w.StartsAt) {
			windows = append(windows, w)
		}
	}

	return windows, nil
}

func deserialize(record []string) (entity.Window, error) {
	var (
		w   entity.Window
		err error
	)

	if w.StartsAt, err = time.Parse(timeFormat, record[0]); err != nil {
		return w, err
	}

	if record[1] == "" {
		return w, nil
	}

	w.EndsAt, err = time.Parse(timeFormat, record[1])
	return w, err
}
