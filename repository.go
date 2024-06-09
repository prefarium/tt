package main

import (
	"encoding/csv"
	"io"
	"os"
	"time"
)

const timeFormat = time.RFC3339

func write(w window) error {
	file, err := os.OpenFile(csvPath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err = writer.Write(serialize(w)); err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}

func read(from, to time.Time) ([]window, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var windows []window
	reader := csv.NewReader(file)
	reader.ReuseRecord = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		w := deserialize(record)
		if !w[0].Before(from) && to.After(w[0]) {
			windows = append(windows, w)
		}
	}

	return windows, nil
}

func serialize(w window) []string {
	return []string{w[0].Format(timeFormat), w[1].Format(timeFormat)}
}

func deserialize(r []string) window {
	return window{mustParse(r[0]), mustParse(r[1])}
}

func mustParse(s string) time.Time {
	t, err := time.Parse(timeFormat, s)
	if err != nil {
		panic(err)
	}

	return t
}
