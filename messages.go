package main

import (
	"bufio"
	"log"
	"os"
)

func Messages(c *AppConfig) <-chan string {
	out := make(chan string)

	go func() {
		file, err := os.Open(c.filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		for {
			file.Seek(0, 0)
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				out <- scanner.Text()
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	return out
}
