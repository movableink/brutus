package main

import (
	"sync"

	"github.com/codegangsta/cli"
)

type AppConfig struct {
	reqsPerSec       int
	threads          int
	numberOfRequests int
	filename         string
	waitGroup        sync.WaitGroup
}

func (ac *AppConfig) Parse(args []string, p DataPusher) error {
	app := cli.NewApp()

	app.Name = "brutus"
	app.Usage = "stab your services in the back...by flooding them with traffic!"

	redisProvider := &RedisProvider{}
	rabbitProvider := &RabbitProvider{}
	nsqProvider := &NSQProvider{}

	app.Commands = []cli.Command{
		redisProvider.Command(ac, p),
		rabbitProvider.Command(ac, p),
		nsqProvider.Command(ac, p),
	}

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "requests, r",
			Value: 200,
			Usage: "limit number of messages sent per second",
		},
		cli.IntFlag{
			Name:  "threads, t",
			Value: 1,
			Usage: "number of concurrently pushing threads",
		},
		cli.StringFlag{
			Name:  "messages, m",
			Value: "messages.json",
			Usage: "location of the messages.json file",
		},
	}

	app.Before = func(c *cli.Context) error {
		ac.reqsPerSec = c.Int("requests")
		ac.threads = c.Int("threads")
		ac.filename = c.String("messages")
		ac.numberOfRequests = 0
		return nil
	}

	return app.Run(args)
}
