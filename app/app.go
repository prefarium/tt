package app

import (
	"path/filepath"
	"time"
	"tt/repository"
	"tt/utils"
)

const fileName = "db.csv"

type App struct {
	repo   *repository.Repo
	offset time.Duration
}

func NewApp(workingDir string, offset int) *App {
	return &App{
		repository.NewRepo(filepath.Join(workingDir, fileName)),
		time.Hour * time.Duration(offset),
	}
}

func (a *App) WorkedToday() (time.Duration, error) {
	now := time.Now()
	todayByOffset := now.Add(-a.offset)
	from := utils.BeginningOfDay(todayByOffset).Add(a.offset)
	to := utils.NextDay(todayByOffset).Add(a.offset)

	windows, err := a.repo.Read(from, to)
	if err != nil {
		return 0, err
	}

	var worked time.Duration
	for _, w := range windows {
		worked += w.EndsAt.Sub(w.StartsAt)
	}

	return worked, nil
}

func (a *App) OpenWindow() error {
	return a.repo.OpenWindow(time.Now())
}

func (a *App) CloseWindow() error {
	return a.repo.CloseWindow(time.Now())
}
