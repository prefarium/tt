package repository

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
	"tt/entity"
)

func TestRepo_OpenWindow(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		content string
	}{
		{
			name: "open window",
			args: args{p("2024-01-01T12:00:00Z")},
			want: fmt.Sprint(
				"2024-01-01T01:00:00Z,2024-01-01T10:00:00Z\n",
				"2024-01-01T12:00:00Z,\n",
			),
			content: "2024-01-01T01:00:00Z,2024-01-01T10:00:00Z\n",
		},
		{
			name: "ignore other opened windows",
			args: args{p("2024-01-01T12:00:00Z")},
			want: fmt.Sprint(
				"2024-01-01T01:00:00Z,\n",
				"2024-01-01T12:00:00Z,\n",
			),
			content: "2024-01-01T01:00:00Z,\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := createTempFile()
			defer os.Remove(file.Name())
			file.WriteString(tt.content)

			r := &Repo{filePath: file.Name()}
			err := r.OpenWindow(tt.args.t)

			if (err != nil) != tt.wantErr {
				t.Errorf("OpenWindow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := os.ReadFile(file.Name())
			if err != nil {
				panic(err)
			}

			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("OpenWindow()\ngot  = %#v\nwant = %#v", string(got), tt.want)
			}
		})
	}
}

func TestRepo_CloseWindow(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		content string
	}{
		{
			name: "close window",
			args: args{p("2024-01-01T15:00:00Z")},
			want: fmt.Sprint(
				"2024-01-01T01:00:00Z,2024-01-01T10:00:00Z\n",
				"2024-01-01T12:00:00Z,2024-01-01T15:00:00Z\n",
			),
			content: fmt.Sprint(
				"2024-01-01T01:00:00Z,2024-01-01T10:00:00Z\n",
				"2024-01-01T12:00:00Z,\n",
			),
		},
		{
			name:    "fail if window is already closed",
			args:    args{p("2024-01-01T15:00:00+03:00")},
			wantErr: true,
			content: fmt.Sprint(
				"2024-01-01T01:00:00+03:00,2024-01-01T10:00:00+03:00\n",
				"2024-01-01T12:00:00+03:00,2024-01-01T14:00:00+03:00\n",
			),
		},
		{
			name:    "fail if no windows present",
			args:    args{p("2024-01-01T15:00:00+03:00")},
			wantErr: true,
			content: "",
		},
		{
			name: "ignore other opened window",
			args: args{p("2024-01-01T15:00:00Z")},
			want: fmt.Sprint(
				"2024-01-01T01:00:00Z,\n",
				"2024-01-01T12:00:00Z,2024-01-01T15:00:00Z\n",
			),
			content: fmt.Sprint(
				"2024-01-01T01:00:00Z,\n",
				"2024-01-01T12:00:00Z,\n",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := createTempFile()
			defer os.Remove(file.Name())
			file.WriteString(tt.content)

			r := &Repo{filePath: file.Name()}
			err := r.CloseWindow(tt.args.t)

			if (err != nil) != tt.wantErr {
				t.Errorf("CloseWindow() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}

			got, err := os.ReadFile(file.Name())
			if err != nil {
				panic(err)
			}

			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("CloseWindow()\ngot  = %#v\nwant = %#v", string(got), tt.want)
			}
		})
	}
}

func TestRepo_Read(t *testing.T) {
	type args struct {
		from time.Time
		to   time.Time
	}

	var (
		from = p("2024-01-02T00:00:00Z")
		to   = p("2024-01-03T00:00:00Z")
	)

	tests := []struct {
		name    string
		args    args
		want    []entity.Window
		wantErr bool
		content string
	}{
		{
			name: "read set range",
			args: args{from, to},
			want: []entity.Window{
				{p("2024-01-02T01:00:00Z"), p("2024-01-02T02:00:00Z")},
				{p("2024-01-02T15:30:30Z"), p("2024-01-02T18:59:59Z")},
			},
			content: fmt.Sprint(
				"2024-01-01T02:00:00Z,2024-01-01T10:00:00Z\n",
				"2024-01-02T01:00:00Z,2024-01-02T02:00:00Z\n",
				"2024-01-02T15:30:30Z,2024-01-02T18:59:59Z\n",
				"2024-01-03T03:00:00Z,2024-01-03T05:00:00Z\n",
			),
		},
		{
			name:    "read windows starting in range and ending out of it",
			args:    args{from, to},
			want:    []entity.Window{{p("2024-01-02T23:59:59Z"), p("2024-01-03T04:00:00Z")}},
			content: "2024-01-02T23:59:59Z,2024-01-03T04:00:00Z\n",
		},
		{
			name:    "do not read windows starting out of range and ending in it",
			args:    args{from, to},
			content: "2024-01-01T23:59:59Z,2024-01-02T04:00:00Z\n",
		},
		{
			name:    "read from range start",
			args:    args{from, to},
			want:    []entity.Window{{p("2024-01-02T00:00:00Z"), p("2024-01-02T10:00:00Z")}},
			content: "2024-01-02T00:00:00Z,2024-01-02T10:00:00Z\n",
		},
		{
			name:    "do not read range end",
			args:    args{from, to},
			content: "2024-01-03T00:00:00Z,2024-01-04T10:00:00Z\n",
		},
		{
			name: "respect time zone",
			args: args{from, to},
			want: []entity.Window{
				{p("2024-01-01T20:00:00-05:00"), p("2024-01-01T21:00:00-05:00")},
				{p("2024-01-03T01:30:30+10:00"), p("2024-01-03T04:59:59+10:00")},
			},
			content: fmt.Sprint(
				"2024-01-01T02:00:00Z,2024-01-01T10:00:00Z\n",
				"2024-01-01T20:00:00-05:00,2024-01-01T21:00:00-05:00\n",
				"2024-01-03T01:30:30+10:00,2024-01-03T04:59:59+10:00\n",
				"2024-01-02T23:00:00-04:00,2024-01-03T01:00:00-04:00\n",
			),
		},
		{
			name: "read not closed windows",
			args: args{from, to},
			want: []entity.Window{
				{p("2024-01-02T01:00:00Z"), time.Time{}},
				{p("2024-01-02T15:30:30Z"), p("2024-01-02T18:59:59Z")},
			},
			content: fmt.Sprint(
				"2024-01-01T02:00:00Z,2024-01-01T10:00:00Z\n",
				"2024-01-02T01:00:00Z,\n",
				"2024-01-02T15:30:30Z,2024-01-02T18:59:59Z\n",
				"2024-01-03T03:00:00Z,2024-01-03T05:00:00Z\n",
			),
		},
		{
			name:    "inadequate file content",
			args:    args{from, to},
			wantErr: true,
			content: "foobar",
		},
		{
			name:    "incorrect date format",
			args:    args{from, to},
			wantErr: true,
			content: "02 Jan 06 15:04 MST,03 Jan 06 15:04 MST",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := createTempFile()
			defer os.Remove(file.Name())
			file.WriteString(tt.content)

			r := &Repo{filePath: file.Name()}
			got, err := r.Read(tt.args.from, tt.args.to)

			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func createTempFile() *os.File {
	tempFile, err := os.CreateTemp("", "test_file.csv")
	if err != nil {
		panic(err)
	}
	return tempFile
}

func p(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
