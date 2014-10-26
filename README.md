gomkr
=====

```gomkr``` is a CLI tool in Go for [mackerel.io](https://mackerel.io).

## Installation

```bash
curl -sL github.com/mackerelio/gomkr/releases/download/latest/gomkr-linux-amd64 > ~/bin/gomkr
```

## Usage

```bash
export MACKEREL_APIKEY=<Put your API key>
```

### On your local host

```bash
gomkr status <hostId> [--json]
gomkr hosts [--name[-n] <hostName>] [--service[-s] <serviceName>] [--role[-r] <roleName>] [--json]
```

```bash
gomkr create [--status[-s] <statusName>] [ [--fullRoleName[-r] <serviceName>:<roleName>] ... ] <hostName>
gomkr update [--status[-s] <statusName>] [ [--fullRoleName[-r] <serviceName>:<roleName>] ... ] <hostId>
```

```bash
gomkr throw --host[-h] <hostId> name\ttime\tvalue [ [name\ttime\tvalue] ... ]
gomkr throw --service[-s] <serviceName> name\ttime\tvalue [ [name\ttime\tvalue] ... ]
```

```bash
gomkr fetch --name[-n] <metricName> [--json] <hostId>
```

```bash
gomkr retire <hostId>
```

### On host running mackerel-agent

You can omit specifing ```<hostId>``` and ```MACKEREL_APIKEY```.
```gomkr``` refers ```/var/lib/mackerel-agent/id``` and ```/etc/mackerel-agent/mackerel-agent.conf``` instead of specifing ```<hostId>```.

```bash
gomkr status
```

```bash
gomkr update [--status[-s] <statusName>] [ [--fullRoleName[-r] <serviceName>:<roleName>] ... ]
```

```bash
gomkr fetch --name[-n] <metricName>
```

```bash
gomkr throw name\ttime\tvalue [ [name\ttime\tvalue] ... ]
```

```bash
gomkr retire
```

## Advanced Usage

```bash
gomkr host list -s My-Service -r proxy -o "id" | xargs -I{} gomkr host status -h {} working
```

```bash
gomkr host list -s My-Service -r proxy -o "id" | xargs -I{} gomkr metric show -h {}
```

## Contribution

1. Fork ([https://github.com/mackerelio/gomkr/fork](https://github.com/mackerelio/gomkr/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request
