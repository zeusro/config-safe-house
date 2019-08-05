
all: build

ARCH?=amd64
VERSION?=v1.0.0
GIT_COMMIT:=$(shell git rev-parse --short HEAD)
IMAGE?=zeusro/config-safe-house:1.0.0

build:
	GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "$(KUBE_EVENTER_LDFLAGS)" -o config-safe-house  github.com/zeusro/config-safe-house

clean:
	rm -f config-safe-house

docker:
	docker build --pull -t $(IMAGE) -f Dockerfile .
	docker push $(IMAGE)