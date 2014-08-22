package main

import (
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

type RabbitProvider struct {
	connection                   *amqp.Connection
	channel                      *amqp.Channel
	server, exchange, routingKey string
}

func (p *RabbitProvider) Command(ac *AppConfig, pusher DataPusher) *cobra.Command {
	command := cobra.Command{
		Use:   "rabbit",
		Short: "push data to a RabbitMQ exchange",
		Run: func(cmd *cobra.Command, args []string) {
			pusher(p, ac)
		},
	}

	command.Flags().StringVarP(&p.server, "server", "s", "amqp://localhost", "rabbit server to push to")
	command.Flags().StringVarP(&p.exchange, "exchange", "e", "brutus", "exchange to use when publishing")
	command.Flags().StringVarP(&p.routingKey, "key", "k", "brutus.messages", "routing key to publish messages to")

	return &command
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
