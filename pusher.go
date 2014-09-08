package main

import (
	"log"
	"time"
)

type DataPusher func(Provider, *AppConfig)

func pusher(p Provider, c *AppConfig) {
	c.waitGroup.Add(c.threads)
	err := p.Connect()
	if err != nil {
		log.Fatal(err)
	}

	messages := Messages(c)
	for i := 0; i < c.threads; i++ {
		go func(p Publisher) {
			var throttle <-chan time.Time
			if c.reqsPerSec > 0 {
				throttle = time.Tick(time.Second / time.Duration(c.reqsPerSec))
			}

			for {
				if c.reqsPerSec > 0 {
					<-throttle
				}

				var message string
				for message = <-messages; !c.msgFilter(message); message = <-messages {
				}

				err := p.Publish(message)
				if err != nil {
					log.Fatal(err)
				} else {
					c.numberOfRequests += 1
				}
			}

			c.waitGroup.Done()
		}(p.NewPublisher())
	}
}
