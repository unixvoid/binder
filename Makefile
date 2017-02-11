GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0
IMAGE_NAME=unixvoid/binder
DOCKER_OPTIONS="--no-cache"
GIT_HASH=$(shell git rev-parse HEAD | head -c 10)

run:
	go run \
		binder/binder.go \
		binder/register.go \
		binder/rotate.go \
		binder/bootstrap.go \
		binder/upload.go \
		binder/set_key.go \
		binder/set_file.go \
		binder/get_key.go \
		binder/get_file.go \
		binder/encrypt.go \
		binder/decrypt.go \
		binder/remove.go \
		binder/garbage_collect.go
stat:
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/binder-$(GIT_HASH)-linux-amd64 binder/*.go

dependencies:
	go get github.com/gorilla/mux
	go get github.com/unixvoid/glogger
	go get gopkg.in/gcfg.v1
	go get gopkg.in/redis.v4
	go get golang.org/x/crypto/sha3

docker:
	$(MAKE) stat
	mkdir stage.tmp/
	cp bin/binder stage.tmp/
	cp config.gcfg stage.tmp/
	cp deps/Dockerfile stage.tmp/
	cd stage.tmp/ && \
		sudo docker build $(DOCKER_OPTIONS) -t $(IMAGE_NAME) .
	@echo "$(IMAGE_NAME) built"

compose:
	cd deps/binder-stack/ && \
		sudo docker-compose up

compose-build:
	cd deps/binder-stack/ && \
		sudo docker-compose build

link-volume:
	cd deps/ && \
		./linkvolume.sh

clean:
	rm -rf bin/
	rm -rf stage.tmp/
