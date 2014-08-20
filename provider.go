package main

import (
	"log"

	"github.com/bitly/go-nsq"
	"github.com/garyburd/redigo/redis"
)

type Provider interface {
	Connect(host string) error
	Publish(topic, msg string) error
}

type RedisProvider struct {
	conn redis.Conn
}

type NSQProvider struct {
	producer *nsq.Producer
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
