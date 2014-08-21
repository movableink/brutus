package main

import (
	"log"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/garyburd/redigo/redis"
	"github.com/streadway/amqp"
)

type Provider interface {
	Connect(host string) error
	Publish(topic, msg string) error
}

// Redis Provider

type RedisProvider struct {
	conn redis.Conn
}

func (r *RedisProvider) Connect(host string) error {
	var err error
	r.conn, err = redis.Dial("tcp", host)

	return err
}

func (r *RedisProvider) Publish(topic, msg string) error {
	_, err := r.conn.Do("rpush", topic, msg)
	return err
}

// NSQ Provider

type NSQProvider struct {
	producer *nsq.Producer
}

func (n *NSQProvider) Connect(host string) error {
	var err error
	config := nsq.NewConfig()
	err = config.Validate()
	if err != nil {
		log.Fatal(err)
	}

	n.producer, err = nsq.NewProducer(host, config)

	return err
}

func (n *NSQProvider) Publish(topic, msg string) error {
	return n.producer.Publish(topic, []byte(msg))
}

// Rabbit Provider

type RabbitProvider struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func (p *RabbitProvider) Connect(host string) error {
	var err error
	p.connection, err = amqp.Dial(host)
	if err != nil {
		log.Fatal(err)
	}

	p.channel, err = p.connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = p.channel.ExchangeDeclare("ojos", "topic", true, true, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (p *RabbitProvider) Publish(topic, msg string) error {
	amqpMsg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         []byte(msg),
	}

	return p.channel.Publish("ojos", topic, true, false, amqpMsg)
}
