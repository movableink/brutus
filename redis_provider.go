package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/cobra"
)

type RedisProvider struct {
	server, list string
	conn         redis.Conn
}

func (p *RedisProvider) Command(ac *AppConfig, pusher DataPusher) *cobra.Command {
	command := cobra.Command{
		Use:   "redis",
		Short: "push data to a redis list",
		Run: func(cmd *cobra.Command, args []string) {
			pusher(p, ac)
		},
	}

	command.Flags().StringVarP(&p.server, "server", "s", "localhost:6379", "redis server to push to")
	command.Flags().StringVarP(&p.list, "list", "l", "brutus_list", "list key to push to")

	return &command
}

func (p *RedisProvider) Connect() error {
	var err error
	p.conn, err = redis.Dial("tcp", p.server)
	return err
}

func (p *RedisProvider) Publish(msg string) error {
	_, err := p.conn.Do("rpush", p.list, msg)
	return err
}
