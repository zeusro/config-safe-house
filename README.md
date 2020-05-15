# config-safe-house

Backup `consul` kv,`kubernetes` yaml and so no.

## usage

### backup once

- set `TEST_CONSUL_Host`

- run test in `sinks/consul/consul_backup_test.go`

- then you can see backup folder in this project

### config

example: [config-default.yaml](/config-default.yaml)

```bash
    nano config.yaml
    docker-compose up
```

docker is just one of the solution,you can also go build and run directly.

## Feature

- [x] Backup consul configs (kv) periodically

## TODO

- [x] 配置解析(2019-08-03)
- [x] 配置替换(2019-08-07 以测试用例形式进行)
- [ ] dingtalk alert

```bash
    export GO111MODULE=on
    go mod vendor
```