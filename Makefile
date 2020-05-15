ARCH        ?=amd64
VERSION     ?=v1.0.0
GIT_COMMIT  :=$(shell git rev-parse --short HEAD)
IMAGE  		?=zeusro/config-safe-house:1.0.1
now         :=$(shell date)

auto_commit:
	git add .
	# 需要注意的是，每行命令在一个单独的shell中执行。这些Shell之间没有继承关系。
	git commit -am "$(now)"
	git pull
	git push

build:
	GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "$(KUBE_EVENTER_LDFLAGS)" -o config-safe-house  github.com/zeusro/config-safe-house

clean:
	rm -f config-safe-house

docker:
	docker build --pull -t $(IMAGE) -f Dockerfile .
	docker push $(IMAGE)

up:
	docker-compose build --no-cache
	docker-compose up --force-recreate 
