package main

import (
	"fmt"
	"log"
	"time"

	"github.com/codegangsta/cli"
	"github.com/streadway/amqp"
)

type RabbitProvider struct {
	connection                   *amqp.Connection
	channel                      *amqp.Channel
	server, exchange, routingKey string
}

func (p *RabbitProvider) Command(ac *AppConfig, pusher DataPusher) cli.Command {
	return cli.Command{
		Name:  "rabbit",
		Usage: "push data to a RabbitMQ exchange",
		Action: func(c *cli.Context) {
			p.server = c.String("server")
			p.exchange = c.String("exchange")
			p.routingKey = c.String("key")

			if p.exchange == "" {
				fmt.Println("ERROR: --exchange, -e argument is required\n")
				cli.ShowCommandHelp(c, "rabbit")
				return
			}

			if p.routingKey == "" {
				fmt.Println("ERROR: --key, -k argument is required\n")
				cli.ShowCommandHelp(c, "rabbit")
				return
			}

			pusher(p, ac)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "server, s",
				Value: "amqp://localhost",
				Usage: "rabbit server to push data to",
			},
			cli.StringFlag{
				Name:  "exchange, e",
				Usage: "exchange to push data to",
			},
			cli.StringFlag{
				Name:  "key, k",
				Usage: "routing key to use when pushing data",
			},
		},
	}
}

func (p *RabbitProvider) Connect() error {
	var err error
	p.connection, err = amqp.Dial(p.server)
	if err != nil {
		log.Fatal(err)
	}

	p.channel, err = p.connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = p.channel.ExchangeDeclare(p.exchange, "topic", true, true, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (p *RabbitProvider) Publish(msg string) error {
	amqpMsg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         []byte(msg),
	}

	return p.channel.Publish(p.exchange, p.routingKey, true, false, amqpMsg)
}
