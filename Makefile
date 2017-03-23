GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0
IMAGE_NAME=unixvoid/binder
DOCKER_OPTIONS="--no-cache"
# standalone dependency locations
ALPINE_FS=https://cryo.unixvoid.com/bin/filesystem/alpine/linux-latest-amd64.rootfs.tar.gz
REDIS_FS=https://cryo.unixvoid.com/bin/redis/filesystem/rootfs.tar.gz
REDIS_BIN=https://cryo.unixvoid.com/bin/redis/3.2.8/redis-server
NGINX_BIN=https://cryo.unixvoid.com/bin/nginx/fancy_index/nginx-1.11.10-linux-amd64

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
		binder/remove.go
stat:
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/binder binder/*.go

dependencies:
	go get github.com/gorilla/mux
	go get github.com/unixvoid/glogger
	go get gopkg.in/gcfg.v1
	go get gopkg.in/redis.v4
	go get golang.org/x/crypto/sha3

prep_api_aci:	stat
	mkdir -p stage.tmp/binder-api-layout/rootfs/
	# copy in binder and its conf
	cp bin/binder stage.tmp/binder-api-layout/rootfs/
	cp config.gcfg stage.tmp/binder-api-layout/rootfs/
	# copy in manifest
	cp deps/manifest.api.json stage.tmp/binder-api-layout/manifest

build_api_aci: prep_api_aci
	# build image
	cd stage.tmp/ && \
		actool build binder-api-layout binder-api.aci && \
		mv binder-api.aci ../
	@echo "binder-api.aci built"

build_travis_api_aci: prep_api_aci
	wget https://github.com/appc/spec/releases/download/v0.8.7/appc-v0.8.7.tar.gz
	tar -zxf appc-v0.8.7.tar.gz
	# build image
	cd stage.tmp/ && \
		../appc-v0.8.7/actool build binder-api-layout binder-api.aci && \
		mv binder-api.aci ../
	rm -rf appc-v0.8.7*
	@echo "binder-api.aci built"

docker:
	$(MAKE) stat
	mkdir stage.tmp/
	cp bin/binder stage.tmp/
	cp config.gcfg stage.tmp/
	cp deps/Dockerfile stage.tmp/
	cd stage.tmp/ && \
		sudo docker build $(DOCKER_OPTIONS) -t $(IMAGE_NAME) .
	@echo "$(IMAGE_NAME) built"

clean:
	rm -rf bin/
	rm -f binder.aci
	rm -rf stage.tmp/
