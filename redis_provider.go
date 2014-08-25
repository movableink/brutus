package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/cobra"
)

type RedisProvider struct {
	server, list string
	pool         *redis.Pool
	conn         redis.Conn
	messageCount int
}

func (p *RedisProvider) Command(ac *AppConfig, pusher DataPusher) *cobra.Command {
	command := cobra.Command{
		Use:   "redis",
		Short: "push data to a redis list",
		Run: func(cmd *cobra.Command, args []string) {
			p.pool = newPool(p.server)
			pusher(p, ac)
		},
	}

	command.Flags().StringVarP(&p.server, "server", "s", "localhost:6379", "redis server to push to")
	command.Flags().StringVarP(&p.list, "list", "l", "brutus_list", "list key to push to")

	return &command
}

func (p *RedisProvider) Connect() error {
	var err error
	p.messageCount = 0
	return err
}

func (p *RedisProvider) Publish(msg string) error {
	conn := p.pool.Get()
	defer conn.Close()
	conn.Send("rpush", p.list, msg)

	p.messageCount++
	if p.messageCount >= 10000 {
		p.messageCount = 0
		conn.Flush()
		_, err := conn.Receive()
		return err
	}
	return nil
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", server)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
