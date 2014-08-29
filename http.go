package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type HTTPProvider struct {
	server      string
	matchTiming bool
}

func (p *HTTPProvider) Command(ac *AppConfig, pusher DataPusher) *cobra.Command {
	command := cobra.Command{
		Use:   "http",
		Short: "replay HTTP requests from a log file",
		Run: func(cmd *cobra.Command, args []string) {
			if p.matchTiming {
				ac.reqsPerSec = 0
			}
			pusher(p, ac)
		},
	}

	command.Flags().StringVarP(&p.server, "server", "s", "http://localhost", "web server to send HTTP requests to")
	command.Flags().BoolVarP(&p.matchTiming, "match-timing", "m", false, "attempt to match the timing of the original messages")

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
	provider     *HTTPProvider
	client       *http.Client
	firstMsgTime time.Time
	startTime    time.Time
	lastDelta    time.Duration
}

type LogMessage struct {
	Message, Timestamp string
}

func recordStart(p *HTTPPublisher, timestamp time.Time) {
	p.startTime = time.Now()
	p.firstMsgTime = timestamp
}

func (p *HTTPPublisher) Publish(msg string) error {
	logMsg := LogMessage{}
	err := json.Unmarshal([]byte(msg), &logMsg)
	if err != nil {
		return err
	}

	if p.provider.matchTiming {
		timestamp, err := time.Parse(time.RFC3339, logMsg.Timestamp)
		if err != nil {
			return err
		}

		// if this is the first message we've processed, record current time and log start time
		// so we can calculate relative timings
		if p.firstMsgTime.IsZero() {
			recordStart(p, timestamp)
		}

		// determine whether the input file has just been restarted
		// if so, reset starting times so relativity between real time and file time is preserved
		currDelta := timestamp.Sub(p.firstMsgTime)
		if currDelta < p.lastDelta {
			recordStart(p, timestamp)
		}
		p.lastDelta = currDelta

		// how long has it been since we began publishing?
		elapsed := time.Now().Sub(p.startTime)

		// sleep until it's time to send the next message (unless it's already time)
		if elapsed < currDelta {
			time.Sleep(currDelta - elapsed)
		}
	}

	words := strings.Fields(logMsg.Message)
	for _, w := range words {
		if strings.HasPrefix(w, "/") {
			resp, err := p.client.Get(p.provider.server + w)

			if resp == nil && err != nil {
				return err
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
