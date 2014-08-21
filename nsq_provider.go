package main

import (
	"fmt"
	"log"

	"github.com/bitly/go-nsq"
	"github.com/codegangsta/cli"
)

type NSQProvider struct {
	producer      *nsq.Producer
	server, topic string
	config        *AppConfig
}

func (p *NSQProvider) Command(ac *AppConfig, pusher DataPusher) cli.Command {
	p.config = ac

	return cli.Command{
		Name:  "nsq",
		Usage: "push data to an NSQ server",
		Action: func(c *cli.Context) {
			p.server = c.String("server")
			p.topic = c.String("topic")

			if p.topic == "" {
				fmt.Println("ERROR: --exchange, -e argument is required\n")
				cli.ShowCommandHelp(c, "rabbit")
				return
			}

			pusher(p, ac)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "server, s",
				Value: "localhost:4150",
				Usage: "NSQ server to push data to",
			},
			cli.StringFlag{
				Name:  "topic, t",
				Usage: "topic to publish data on",
			},
		},
	}
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
