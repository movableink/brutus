package main

import "github.com/spf13/cobra"

type Provider interface {
	Command(*AppConfig, DataPusher) *cobra.Command
	Connect() error
	NewPublisher() Publisher
}
