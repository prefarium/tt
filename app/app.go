package app

import (
	"time"
	"tt/entity"
	"tt/repository"
	"tt/utils"
)

type App struct {
	repo     repository.Repo
	offset   time.Duration
	location *time.Location
}

func NewApp(workingDir string, offset time.Duration, location *time.Location) App {
	return App{
		repository.NewRepo(workingDir),
		offset,
		location,
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

func (a *App) WorkedThisWeek() (entity.Week, error) {
	today := utils.BeginningOfDay(time.Now().In(a.location))
	monday := today.AddDate(0, 0, -dayNumber(today))
	week := entity.Week{}

	for i := range week {
		from := monday.AddDate(0, 0, i)
		to := monday.AddDate(0, 0, i+1)

		windows, err := a.repo.Read(from, to)
		if err != nil {
			return week, err
		}

		week[i] = entity.Day(sumWindows(windows))
	}

	return week, nil
}

func dayNumber(t time.Time) int {
	wd := t.Weekday()
	switch wd {
	case time.Sunday:
		return 6
	default:
		return int(wd) - 1
	}
}

func sumWindows(w []entity.Window) time.Duration {
	var sum time.Duration
	for i := range w {
		sum += w[i].EndsAt.Sub(w[i].StartsAt)
	}
	return sum
}
