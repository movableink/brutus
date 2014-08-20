Usage:
```
  ./ [OPTIONS] [TOPIC]
  -f="": Source Data Filename (required)
  -h=":6379": hostname
  -t=1: Number of threads
  -c=: Number of concurrent requests/sec (per thread)
  -p=: Provider Name
```

Example:
```
$ ./brutus -t 4 -c 200 -f messages.json -p "nsq" -h some-remote-server.com:4567 my_topic
```

Development:
1. Install https://github.com/tools/godep
1. run `godep get`
