GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0

run:
	go run \
		binder/binder.go \
		binder/register.go \
		binder/upload.go \
		binder/remove.go
