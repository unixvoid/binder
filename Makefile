GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0
IMAGE_NAME=unixvoid/binder
DOCKER_OPTIONS="--no-cache"

run:
	go run \
		binder/binder.go \
		binder/register.go \
		binder/rotate.go \
		binder/bootstrap.go \
		binder/upload.go \
		binder/remove.go
stat:
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/binder binder/*.go

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
