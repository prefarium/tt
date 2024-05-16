package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tt/app"
	"tt/utils"
)

func main() {
	now := time.Now()
	a := app.NewApp("/Users/anton/GolandProjects/tt/tmp", 3)

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
