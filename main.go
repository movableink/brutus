package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	var threadCt, numberOfRequests, reqsSec int
	var host, filename, providerName string

	flag.IntVar(&threadCt, "t", 1, "Number of threads")
	flag.IntVar(&reqsSec, "r", 200, "Requests a second")
	flag.StringVar(&host, "h", ":6379", "Redis host")
	flag.StringVar(&filename, "f", "", "Source Data Filename (required)")
	flag.StringVar(&providerName, "p", "redis", "Load provider (redis, nsq)")
	flag.Parse()
	topic := flag.Arg(0)

	if len(topic) == 0 {
		log.Fatalln("[TOPIC] is required")
		return
	}
	if len(filename) == 0 {
		log.Fatal(filename, "-f [filename] is required.")
		return
	}

	var wg sync.WaitGroup
	wg.Add(threadCt)

	log.Println("Num Threads: ", threadCt)

	for i := 0; i < threadCt; i++ {
		go func() {
			throttle := time.Tick(time.Second / time.Duration(reqsSec))

			var provider Provider
			switch providerName {
			case "nsq":
				provider = &NSQProvider{}
			default:
				provider = &RedisProvider{}
			}

			provider.Connect(host)

			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			for {
				file.Seek(0, 0)
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					<-throttle
					err := provider.Publish(topic, scanner.Text())

					if err != nil {
						log.Fatal(err)
					} else {
						numberOfRequests += 1
					}
				}

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
			}

			wg.Done()
		}()
	}

	printTick := time.Tick(time.Second)
	go func() {
		for {
			<-printTick
			log.Println("Req/s", numberOfRequests)
			numberOfRequests = 0
		}
	}()

	wg.Wait()
}
