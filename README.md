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
gomkr status <hostId>
```

```bash
gomkr hosts --service My-Service --role db-master --role db-slave
```

```bash
gomkr create --status working -R My-Service:db-master mydb001
gomkr update --status maintenance --role My-Service:db-master <hostId>
```

```bash
cat <<EOF | gomkr throw --host <hostId>
<name>  <time>  <value>
<name>  <time>  <value>
EOF
...

cat <<EOF | gomkr throw --service My-Service
<name>  <time>  <value>
<name>  <time>  <value>
EOF
...
```

```bash
gomkr fetch --name loadavg5 <hostId>
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
gomkr hosts -s My-Service -r proxy | jq -r '.[].id'  | xargs -I{} gomkr update --st working {}
```

## Contribution

1. Fork ([https://github.com/mackerelio/gomkr/fork](https://github.com/mackerelio/gomkr/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request
