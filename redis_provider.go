package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/garyburd/redigo/redis"
)

type RedisProvider struct {
	server, list string
	conn         redis.Conn
}

func (p *RedisProvider) Command(ac *AppConfig, pusher DataPusher) cli.Command {
	return cli.Command{
		Name:  "redis",
		Usage: "push data to a redis list",
		Action: func(c *cli.Context) {
			p.server = c.String("server")
			p.list = c.String("list")

			if p.list == "" {
				fmt.Println("ERROR: --list, -l argument is required\n")
				cli.ShowCommandHelp(c, "redis")
				return
			}

			pusher(p, ac)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "server, s",
				Value: "localhost:6379",
				Usage: "redis server to push data to",
			},
			cli.StringFlag{
				Name:  "list, l",
				Usage: "list key in redis to push data to",
			},
		},
	}
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
