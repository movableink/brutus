package main

import "github.com/codegangsta/cli"

type Provider interface {
	Command(*AppConfig, DataPusher) cli.Command
	Connect() error
	Publish(msg string) error
}
