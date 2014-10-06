#### Development
1. Install https://github.com/tools/godep
1. `go get github.com/movableink/brutus`

_... in $GOPATH/src/github.com/movableink/brutus_

1. run `godep get`

### Build it Locally

`go install` - this will place brutus in your `$GOBIN` directory

### Usage

Run `brutus help` for a list of commands

```
Usage:
  brutus [command]

Available Commands:
  redis                     push data to a redis list
  rabbit                    push data to a RabbitMQ exchange
  nsq                       push data to an NSQ daemon
  kafka                     push data to a Kafka topic
  http                      replay HTTP requests from a log file
  help [command]            Help about any command

 Available Flags:
  -c, --concurrency=1: number of pusher threads to create
  -f, --filename="messages.json": file containing message data
      --help=false: help for brutus
  -r, --requests=200: target number of req/s (per thread)

Use "brutus help [command]" for more information about that command.
```

Run "brutus help [command]" for more information about that command.


### Build it for linux/amd64
We use [gox](https://github.com/mitchellh/gox) as a cross-compilation tool.

1. `gox -osarch="linux/amd64"`
