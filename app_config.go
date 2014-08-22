package main

import (
	"sync"

	"github.com/spf13/cobra"
)

type AppConfig struct {
	reqsPerSec       int
	threads          int
	numberOfRequests int
	filename         string
	waitGroup        sync.WaitGroup
}

func (ac *AppConfig) RunCli(p DataPusher) error {
	var command = &cobra.Command{
		Use:   "brutus",
		Short: "stab your services in the back...by flooding them with traffic!",
	}

	command.PersistentFlags().StringVarP(&ac.filename, "filename", "f", "messages.json", "file containing message data")
	command.PersistentFlags().IntVarP(&ac.threads, "concurrency", "c", 1, "number of pusher threads to create")
	command.PersistentFlags().IntVarP(&ac.reqsPerSec, "requests", "r", 200, "target number of req/s (per thread)")

	redisProvider := &RedisProvider{}
	rabbitProvider := &RabbitProvider{}
	nsqProvider := &NSQProvider{}

	command.AddCommand(redisProvider.Command(ac, p))
	command.AddCommand(rabbitProvider.Command(ac, p))
	command.AddCommand(nsqProvider.Command(ac, p))

	return command.Execute()
}
