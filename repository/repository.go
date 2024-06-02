package repository

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"time"
	"tt/entity"
)

const timeFormat = time.RFC3339

type Repo struct {
	filePath string
}

func NewRepo(filePath string) Repo {
	return Repo{filePath: filePath}
}

func (r *Repo) OpenWindow(t time.Time) error {
	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err = writer.Write(serialize(t, time.Time{})); err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}

func (r *Repo) CloseWindow(t time.Time) error {
	file, err := os.OpenFile(r.filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	lastIdx := len(records) - 1
	if lastIdx < 0 {
		return errors.New("no windows found")
	}

	window, err := deserialize(records[lastIdx])
	if err != nil {
		return err
	} else if window.IsClosed() {
		return errors.New("window is already closed")
	}

	records[lastIdx] = serialize(window.StartsAt, t)
	writer := csv.NewWriter(file)
	file.Seek(0, io.SeekStart)
	for _, record := range records {
		if err = writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
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

func serialize(from, to time.Time) []string {
	var window []string
	window = append(window, from.Format(timeFormat))

	if to.IsZero() {
		window = append(window, "")
	} else {
		window = append(window, to.Format(timeFormat))
	}

	return window
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
