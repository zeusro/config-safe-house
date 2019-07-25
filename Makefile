
all: build

ARCH?=amd64
VERSION?=v1.0.0
GIT_COMMIT:=$(shell git rev-parse --short HEAD)

build:
	GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "$(KUBE_EVENTER_LDFLAGS)" -o config-safe-house  github.com/zeusro/config-safe-house

clean:
	rm -f config-safe-house

docker:
	# docker build --pull -t zeusro/config-safe-house:$(VERSION)-$(GIT_COMMIT) -f Dockerfile .
	docker build --pull -t zeusro/config-safe-house:1.0.0 -f Dockerfile .