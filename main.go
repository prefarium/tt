package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tt/app"
	"tt/entity"
	"tt/utils"
)

func main() {
	conf, err := parseConfig()
	if err != nil {
		fmt.Printf("parsing config: %s\n", err)
	}

	a, err := initApp(conf)
	if err != nil {
		fmt.Printf("init app: %s\n", err)
	}

	if len(os.Args) == 1 {
		trackTime(a)
	} else if os.Args[1] == "week" {
		calculateWeek(a)
	}
}

func parseConfig() (entity.Config, error) {
	var c entity.Config

	file, err := os.Open("config.json")
	if err != nil {
		return c, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func initApp(c entity.Config) (app.App, error) {
	if _, err := os.Stat(c.CsvPath); err != nil {
		return app.App{}, err
	}

	offset := time.Minute * time.Duration(c.Offset)

	location, err := time.LoadLocation(c.Location)
	if err != nil {
		return app.App{}, err
	}

	return app.NewApp(c.CsvPath, offset, location), nil
}

func trackTime(a app.App) {
	now := time.Now()

	worked, err := a.WorkedToday()
	if err != nil {
		fmt.Printf("failed to calculate time worked today: %s\n", err)
		return
	}

	if err := a.OpenWindow(); err != nil {
		fmt.Printf("failed to start tracking: %s\n", err)
		return
	}

	ticker := time.NewTicker(time.Second)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			elapsed := time.Since(now) + worked
			fmt.Print(fmt.Sprintf("\r%s", utils.PrettyDuration(elapsed)))
		case <-exit:
			ticker.Stop()
			if err := a.CloseWindow(); err != nil {
				fmt.Printf("failed to finish tracking: %s\n", err)
			}
			return
		}
	}
}

func calculateWeek(a app.App) {
	week, err := a.WorkedThisWeek()
	if err != nil {
		fmt.Printf("failed to calculate week: %v\n", err)
	}

	days := map[int]string{
		0: "Mon",
		1: "Tue",
		2: "Wed",
		3: "Thu",
		4: "Fri",
		5: "Sat",
		6: "San",
	}

	for i, d := range week {
		fmt.Printf("%s %s\n", days[i], d)
	}

	fmt.Printf("---\nSum %s\n", utils.PrettyDuration(week.Total()))
}
