# config-safe-house

Backup `consul` kv,`kubernetes` yaml and so no.

## usage

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
- [ ] 配置替换
- [ ] dingtalk alert

    export GO111MODULE=on
    go mod vendor
