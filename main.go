package main

import (
	"log"
	"os"
	"time"
)

func main() {
	config := &AppConfig{}
	err := config.Parse(os.Args, pusher)
	if err != nil {
		log.Fatal(err)
	}

	printTick := time.Tick(time.Second)
	go func() {
		for {
			<-printTick
			log.Println("Req/s", config.numberOfRequests)
			config.numberOfRequests = 0
		}
	}()

	config.waitGroup.Wait()
}
