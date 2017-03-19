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

prep_aci:	stat
	mkdir -p stage.tmp/
	cp -R deps/binder-layout stage.tmp/
	# alpine fs
	wget -O alpinefs.tar.gz $(ALPINE_FS)
	tar -xzf alpinefs.tar.gz -C stage.tmp/binder-layout/rootfs/
	rm alpinefs.tar.gz
	# redis fs
	wget -O redisfs.tar.gz $(REDIS_FS)
	tar -xzf redisfs.tar.gz -C stage.tmp/binder-layout/rootfs/
	rm redisfs.tar.gz
	# redis bin
	wget -O stage.tmp/binder-layout/rootfs/bin/redis $(REDIS_BIN)
	chmod +x stage.tmp/binder-layout/rootfs/bin/redis
	# nginx bin
	wget -O stage.tmp/binder-layout/rootfs/bin/nginx $(NGINX_BIN)
	chmod +x stage.tmp/binder-layout/rootfs/bin/nginx
	# add binder + misc files
	cp bin/* stage.tmp/binder-layout/rootfs/bin/binder
	cp config.gcfg stage.tmp/binder-layout/rootfs/
	cp deps/redis.conf stage.tmp/binder-layout/rootfs/
	cp deps/run.sh stage.tmp/binder-layout/rootfs/
	#touch stage.tmp/binder-layout/rootfs/nginx/log/error.log
	mkdir -p stage.tmp/binder-layout/rootfs/uploads
	mkdir -p stage.tmp/binder-layout/rootfs/nginx/log/
	mkdir -p stage.tmp/binder-layout/rootfs/nginx/conf/
	mkdir -p stage.tmp/binder-layout/rootfs/nginx/data/
	cp -R deps/nginx/nginx_fancyindex_data/* stage.tmp/binder-layout/rootfs/nginx/data/
	cp deps/nginx/nginx.conf stage.tmp/binder-layout/rootfs/nginx/conf/
	cp deps/nginx/mime.types stage.tmp/binder-layout/rootfs/nginx/conf/
	cp deps/manifest.json stage.tmp/binder-layout/manifest

build_aci: prep_aci
	# build image
	cd stage.tmp/ && \
		actool build binder-layout binder.aci && \
		mv binder.aci ../
	@echo "binder.aci built"

build_travis_aci: prep_aci
	wget https://github.com/appc/spec/releases/download/v0.8.7/appc-v0.8.7.tar.gz
	tar -zxf appc-v0.8.7.tar.gz
	# build image
	cd stage.tmp/ && \
		../appc-v0.8.7/actool build binder-layout binder.aci && \
		mv binder.aci ../
	rm -rf appc-v0.8.7*
	@echo "binder.aci built"

docker:
	$(MAKE) stat
	mkdir stage.tmp/
	cp bin/binder stage.tmp/
	cp config.gcfg stage.tmp/
	cp deps/Dockerfile stage.tmp/
	cd stage.tmp/ && \
		sudo docker build $(DOCKER_OPTIONS) -t $(IMAGE_NAME) .
	@echo "$(IMAGE_NAME) built"

link-volume:
	cd deps/ && \
		./linkvolume.sh

clean:
	rm -rf bin/
	rm -f binder.aci
	rm -rf stage.tmp/
