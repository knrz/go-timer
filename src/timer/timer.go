package main

import (
	"errors"
	"os"
	"runtime"
	"strconv"
	"time"
)

import (
	"github.com/codegangsta/cli"
	"gopkg.in/cheggaaa/pb.v1"
)

func total(c *cli.Context) (t int64) {
	if c.NArg() == 1 {
		ti, err := strconv.Atoi(c.Args()[0])
		if err != nil {
			return int64(-1)
		}
		t = int64(ti)
	} else {
		t += int64(c.Int("seconds"))
		t += int64(c.Int("minutes") * 60)
		t += int64(c.Int("hours") * 3600)
		t += int64(c.Int("days") * 86400)
	}
	return
}

func progressBar(t int64, c *cli.Context) *pb.ProgressBar {
	bar := pb.New64(t)
	bar.ShowPercent = false
	bar.ShowCounters = false
	bar.SetRefreshRate(500)
	bar.ShowSpeed = false
	bar.Format(c.String("format"))
	bar.Start()
	return bar
}

func run(c *cli.Context) error {
	t := total(c)
	if t < 0 {
		return errors.New("Time can't be negative")
	}

	bar := progressBar(t, c)

	ticker := time.NewTicker(time.Second)
	go func() {
		for _ = range ticker.C {
			bar.Increment()
		}
	}()

	time.Sleep(time.Duration(t) * time.Second)
	ticker.Stop()

	var bell string
	if runtime.GOOS != "windows" {
		bell = "\a"
	}

	bar.FinishPrint(bell + c.String("message"))
	return nil
}

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Name = "Timer"
	app.Usage = "A simple timer"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "days",
			Usage: "Number of days",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "hours",
			Usage: "Number of hours",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "minutes",
			Usage: "Number of minutes",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "seconds",
			Usage: "Number of seconds",
			Value: 0,
		},
		cli.StringFlag{
			Name:  "message",
			Usage: "Message to print when the timer's finished",
			Value: "Time's up!",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "Specify the format as a 5-character long string, [start][progress][head][left][finish]",
			Value: "==>  ",
		},
	}
	app.Action = run
	app.Run(os.Args)
}
