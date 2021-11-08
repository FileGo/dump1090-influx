# dump1090-influx

[![Go Report Card](https://goreportcard.com/badge/github.com/FileGo/dump1090-influx)](https://goreportcard.com/report/github.com/FileGo/dump1090-influx) ![build](https://github.com/FileGo/dump1090-influx/workflows/build/badge.svg) ![tests](https://github.com/FileGo/dump1090-influx/workflows/tests/badge.svg) ![docker](https://img.shields.io/docker/pulls/filego/dump1090-influx.svg)

This program periodically retrieves JSON data from [dump1090-mutability](https://github.com/adsbxchange/dump1090-mutability) and stores it in InfluxDB.

The easiest way to use this is through Docker with `docker-compose`:

```yaml
version: "3.6"
services:
    dump1090-influx:
        container_name: dump1090-influx
        build: .
        restart: unless_stopped
        environment:
            - HOST=http://localhost/dump1090/data/stats.json
            - INFLUX_URL=http://localhost:8086
            - INFLUX_TOKEN=
            - INFLUX_ORG=
            - INFLUX_BUCKET=dump1090 # Database name for InfluxDB v1
            - POLL_TIME=10s
```