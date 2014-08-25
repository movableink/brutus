package main

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

type KafkaProvider struct {
	client    *sarama.Client
	producer  *sarama.Producer
	server    string
	zookeeper string
	topic     string
}

func (p *KafkaProvider) Command(ac *AppConfig, pusher DataPusher) *cobra.Command {
	command := cobra.Command{
		Use:   "kafka",
		Short: "push data to a Kafka topic",
		Run: func(cmd *cobra.Command, args []string) {
			pusher(p, ac)
		},
	}

	command.Flags().StringVarP(&p.server, "server", "s", "localhost:9092", "kafka server to push to")
	command.Flags().StringVarP(&p.zookeeper, "zookeeper", "z", "localhost:2181", "zookeeper server to register topic with")
	command.Flags().StringVarP(&p.topic, "topic", "t", "brutus", "topic to publish data one")

	return &command
}

func (p *KafkaProvider) Connect() error {
	var err error
	p.client, err = sarama.NewClient("brutus", []string{p.server}, nil)
	if err != nil {
		log.Fatal(err)
	}

	p.producer, err = sarama.NewProducer(p.client, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		outchan := p.producer.Errors()
		for {
			err := <-outchan
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	return err
}

func (p *KafkaProvider) Publish(msg string) error {
	return p.producer.QueueMessage(p.topic, sarama.StringEncoder("brutus"), sarama.StringEncoder(msg))
}
