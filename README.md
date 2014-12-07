mkr
=====

mkr - A fast Mackerel client in Go.

# DESCRIPTION

mkr is a command-line interface tool for [Mackerel API](http://help-ja.mackerel.io/entry/spec/api/v0) written in Go language.
mkr helps you to free your daily troublesome server operations and accelarates to leverage Mackerel and the Unix tools.
mkr output format is JSON, so you can filter it by JSON processor such as [jq](http://stedolan.github.io/jq/).

# INSTALLATION

  $ go get github.com/y-uuki/mkr
  $ go install github.com/y-uuki/mkr

## TODO

```bash
$ curl -sL github.com/mackerelio/mkr/releases/download/latest/mkr-linux-amd64 > ~/bin/mkr && chmod +x ~/bin/mkr
```

# USAGE

Set MACKEREL_APIKEY environment variable, but you don't have to set MACKEREL_APIKEY on your host running [mackerel-agent](https://github.com/mackerelio/mackerel-agent). For more details, see below.

```bash
export MACKEREL_APIKEY=<Put your API key>
```

## EXAMPLES

```
$ mkr status <hostId>
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
$ mkr hosts --service My-Service --role proxy
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
mkr create --status working -R My-Service:db-master mydb001
mkr update --status maintenance --role My-Service:db-master <hostId>
```

```
cat <<EOF | mkr throw --host <hostId>
<name>  <time>  <value>
<name>  <time>  <value>
EOF
...

cat <<EOF | mkr throw --service My-Service
<name>  <time>  <value>
<name>  <time>  <value>
EOF
...
```

```
mkr fetch --name loadavg5 2eQGDXqtoXs
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
mkr retire <hostId> ...
```

### Examples (on host running mackerel-agent)

You can omit specifing <hostId> and MACKEREL_APIKEY.
mkrrefers /var/lib/mackerel-agent/id and /etc/mackerel-agent/mackerel-agent.conf instead of specifing <hostId>.

```
mkr status
```

```
mkr update --status maintenance <hostIds>...
```

```
mkr fetch -n loadavg5
```

```bash
cat <<EOF | mkr throw --host <hostId>
<name>  <time>  <value>
EOF
```

```
mkr retire
```

## ADVANCED USAGE

```bash
$ mkr update --st working $(mkr hosts -s My-Service -r proxy | jq -r '.[].id')
```

# CONTRIBUTION

1. Fork ([https://github.com/mackerelio/mkr/fork](https://github.com/mackerelio/mkr/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request


License
----------

Copyright 2014 Hatena Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
