gomkr
=====

gomkr - A fast Mackerel client in Go.

# DESCRIPTION

gomkr is a command-line interface tool for [Mackerel API](http://help-ja.mackerel.io/entry/spec/api/v0) written in Go language.
gomkr helps you to free your daily troublesome server operations and accelarates to leverage Mackerel and the Unix tools.
gomkr output format is JSON, so you can filter it by JSON processor such as [jq](http://stedolan.github.io/jq/).

## Installation

```bash
$ curl -sL github.com/mackerelio/gomkr/releases/download/latest/gomkr-linux-amd64 > ~/bin/gomkr && chmod +x ~/bin/gomkr
```

## Usage

Set MACKEREL_APIKEY environment variable, but you don't have to set MACKEREL_APIKEY on your host running [mackerel-agent](https://github.com/mackerelio/mackerel-agent). For more details, see below.

```bash
export MACKEREL_APIKEY=<Put your API key>
```

### Examples

```
$ gomkr status <hostId>
{
    "id": "2eQGEaLxiYU",
    "name": "myproxy001",
    "status": "standby",
    "roleFullnames": [
        "My-Service:proxy"
    ],
    "isRetired": false,
    "createdAt": "Nov 15, 2014 at 9:41pm (JST)"
}
```

```
$ gomkr hosts --service My-Service --role proxy
[
    {
        "id": "2eQGEaLxiYU",
        "name": "myproxy001",
        "status": "standby",
        "roleFullnames": [
            "My-Service:proxy"
        ],
        "isRetired": false,
        "createdAt": "Nov 15, 2014 at 9:41pm (JST)"
    },
    {
        "id": "2eQGDXqtoXs",
        "name": "myproxy002",
        "status": "standby",
        "roleFullnames": [
            "My-Serviceg:proxy"
        ],
        "isRetired": false,
        "createdAt": "Nov 15, 2014 at 9:41pm (JST)"
    },
]
```

```
gomkr create --status working -R My-Service:db-master mydb001
gomkr update --status maintenance --role My-Service:db-master <hostId>
```

```
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

```
gomkr fetch --name loadavg5 2eQGDXqtoXs
{
    "2eQGDXqtoXs": {
        "loadavg5": {
            "time": 1416061500,
            "value": 0.025
        }
    }
}
```

```
gomkr retire <hostId> ...
```

### Examples (on host running mackerel-agent)

You can omit specifing <hostId> and MACKEREL_APIKEY.
gomkrrefers /var/lib/mackerel-agent/id and /etc/mackerel-agent/mackerel-agent.conf instead of specifing <hostId>.

```
gomkr status
```

```
gomkr update -st maintenance <hostIds>...
```

```
gomkr fetch -n loadavg5
```

```bash
cat <<EOF | gomkr throw --host <hostId>
<name>  <time>  <value>
EOF
```

```
gomkr retire
```

## Advanced Usage

```bash
$ gomkr update --st working $(gomkr hosts -s My-Service -r proxy | jq -r '.[].id')
```

## Contribution

1. Fork ([https://github.com/mackerelio/gomkr/fork](https://github.com/mackerelio/gomkr/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request
