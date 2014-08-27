package main

import (
	"bufio"
	"log"
	"os"
	"time"
)

type DataPusher func(Provider, *AppConfig)

func pusher(p Provider, c *AppConfig) {
	c.waitGroup.Add(c.threads)
	err := p.Connect()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < c.threads; i++ {
		go func(p Publisher) {
			throttle := time.Tick(time.Second / time.Duration(c.reqsPerSec))

			file, err := os.Open(c.filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			for {
				file.Seek(0, 0)
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					<-throttle

					err := p.Publish(scanner.Text())
					if err != nil {
						log.Fatal(err)
					} else {
						c.numberOfRequests += 1
					}
				}

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
			}

			c.waitGroup.Done()
		}(p.NewPublisher())
	}
}
