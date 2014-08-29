package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

type HTTPProvider struct {
	server string
}

func (p *HTTPProvider) Command(ac *AppConfig, pusher DataPusher) *cobra.Command {
	command := cobra.Command{
		Use:   "http",
		Short: "replay HTTP requests from a log file",
		Run: func(cmd *cobra.Command, args []string) {
			pusher(p, ac)
		},
	}

	command.Flags().StringVarP(&p.server, "server", "s", "http://localhost", "web server to send HTTP requests to")

	ac.msgFilter = func(message string) bool {
		return strings.Index(message, "GET") >= 0
	}

	return &command
}

func (p *HTTPProvider) Connect() error {
	return nil
}

func (p *HTTPProvider) NewPublisher() Publisher {
	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("never redirect me")
		},
	}
	return &HTTPPublisher{provider: p, client: c}
}

type HTTPPublisher struct {
	provider *HTTPProvider
	client   *http.Client
}

type LogMessage struct {
	Message, Timestamp string
}

func (p *HTTPPublisher) Publish(msg string) error {
	logMsg := LogMessage{}

	err := json.Unmarshal([]byte(msg), &logMsg)
	if err != nil {
		log.Fatal(err)
	}

	words := strings.Fields(logMsg.Message)
	for _, w := range words {
		if strings.HasPrefix(w, "/") {
			resp, err := p.client.Get(p.provider.server + w)

			if resp == nil && err != nil {
				fmt.Println("ERROR:", err)
			}

			if resp != nil {
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}

			return nil
		}
	}

	return err
}
