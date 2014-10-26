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
gomkr status <hostId> [--json]
gomkr hosts [--name[-n] <hostName>] [--service[-s] <serviceName>] [--role[-r] <roleName>] [--json]
```

```
gomkr create [--status[-s] <statusName>] [ [--fullRoleName[-r] <serviceName>:<roleName>] ... ] <hostName>
gomkr update [--status[-s] <statusName>] [ [--fullRoleName[-r] <serviceName>:<roleName>] ... ] <hostId>
```

```
gomkr throw --host[-h] <hostId> name\ttime\tvalue [ [name\ttime\tvalue] ... ]
gomkr throw --service[-s] <serviceName> name\ttime\tvalue [ [name\ttime\tvalue] ... ]
```

```
gomkr fetch --name[-n] <metricName> [--json] <hostId>
```

```
gomkr retire <hostId>
```

### On host running mackerel-agent

You can omit specifing ```<hostId>``` and ```MACKEREL_APIKEY```.
```gomkr``` refers ```/var/lib/mackerel-agent/id``` and ```/etc/mackerel-agent/mackerel-agent.conf``` instead of specifing ```<hostId>```.

```
gomkr status
```

```
gomkr update [--status[-s] <statusName>] [ [--fullRoleName[-r] <serviceName>:<roleName>] ... ]
```

```
gomkr fetch --name[-n] <metricName>
```

```
gomkr throw name\ttime\tvalue [ [name\ttime\tvalue] ... ]
```

```
gomkr retire
```

## Advanced Usage

```
gomkr host list -s My-Service -r proxy -o "id" | xargs -I{} gomkr host status -h {} working
```

```
gomkr host list -s My-Service -r proxy -o "id" | xargs -I{} gomkr metric show -h {}
```
