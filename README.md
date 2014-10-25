gomkr
=====

```gomkr``` is a CLI tool in Go for [mackerel.io](https://mackerel.io).

## Installation

```
curl -sL github.com/mackerelio/gomkr/releases/download/latest/gomkr-linux-amd64 > ~/bin/gomkr
```

## Usage

```
export MACKEREL_APIKEY=<Put your API key>
```

### On your local host

```
gomkr host show <hostId> [--json]
gomkr host list [--name[-n] <hostName>] [--service[-s] <serviceName>] [--role[-r] <roleName>] [--json]
gomkr host status --host[-h] <hostId> --status[-s] <statusName>
```

```
gomkr metric --host[-h] <hostId> name\ttime\tvalue [ [name\ttime\tvalue] ... ]
gomkr metric --service[-s] <serviceName> name\ttime\tvalue [ [name\ttime\tvalue] ... ]
```

```
gomkr metric show --host[-h] <hostId> [--name[-n] <metricName>] [--json]
```

```
gomkr host retire --host[-h] <hostId>
```

### On host running mackerel-agent

You can omit specifing ```<hostId>``` and ```MACKEREL_APIKEY```.
```gomkr``` refers ```/var/lib/mackerel-agent/id``` and ```/etc/mackerel-agent/mackerel-agent.conf``` instead of specifing ```<hostId>```.

```
gomkr host info
```

```
gomkr host status --status[-s] <statusName>
```

```
gomkr metric name\ttime\tvalue [ [name\ttime\tvalue] ... ]
```

## Advanced Usage

```
gomkr host list -s My-Service -r proxy -o "id" | xargs -I{} gomkr host status -h {} working
```

```
gomkr host list -s My-Service -r proxy -o "id" | xargs -I{} gomkr metric show -h {}
```
