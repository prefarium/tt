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
	a := newApp()

	worked, err := a.WorkedToday()
	if err != nil {
		panic(err)
	}

	if err := a.OpenWindow(); err != nil {
		panic(err)
	}

	for {
		timePassed := time.Now().Add(worked).Sub(now)
		fmt.Print(fmt.Sprintf("\r%s", utils.PrettyDuration(timePassed)))
		time.Sleep(time.Second)
	}
}

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := newApp().CloseWindow(); err != nil {
			panic(err)
		}
		fmt.Println()
		os.Exit(0)
	}()
}

func newApp() *app.App {
	return app.NewApp("/Users/anton/GolandProjects/tt/tmp", 3)
}
