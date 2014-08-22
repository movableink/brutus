package main

import (
	"log"

	"github.com/bitly/go-nsq"
	"github.com/spf13/cobra"
)

type NSQProvider struct {
	producer      *nsq.Producer
	server, topic string
	config        *AppConfig
}

func (p *NSQProvider) Command(ac *AppConfig, pusher DataPusher) *cobra.Command {
	command := cobra.Command{
		Use:   "nsq",
		Short: "push data to an NSQ server",
		Run: func(cmd *cobra.Command, args []string) {
			pusher(p, ac)
		},
	}

	command.Flags().StringVarP(&p.server, "server", "s", "localhost:4150", "NSQ server to push to")
	command.Flags().StringVarP(&p.topic, "topic", "t", "brutus.messages", "topic to publish messages on")

	return &command
}

func (p *NSQProvider) Connect() error {
	config := nsq.NewConfig()
	err := config.Validate()
	if err != nil {
		log.Fatal(err)
	}

	p.producer, err = nsq.NewProducer(p.server, config)

	return err
}

func (p *NSQProvider) Publish(msg string) error {
	return p.producer.Publish(p.topic, []byte(msg))
}
