all: push

USERSPACE?=nikita-vanyasin
GOOS?=linux

APPSERVER?=build/server/go-server
APPDAEMON?=build/daemon/go-daemon

PACKAGESERVER?=github.com/${USERSPACE}/go-web-course/cmd/server
PACKAGEDAEMON?=github.com/${USERSPACE}/go-web-course/cmd/daemon

CURRENTDIR=$(shell pwd)

vendor: clean
	go get -u github.com/golang/dep/cmd/dep \
	&& dep ensure

build: vendor
	CGO_ENABLED=0 GOOS=${GOOS} go build -a -installsuffix cgo -o ${APPSERVER} ${PACKAGESERVER} \
	&& CGO_ENABLED=0 GOOS=${GOOS} go build -a -installsuffix cgo -o ${APPDAEMON} ${PACKAGEDAEMON}

run: build
	CONTENT_FOLDER_PATH=${CURRENTDIR}/content docker-compose up -d

clean:
	rm -f ${APPSERVER} \
	&& rm -f ${APPDAEMON} \
	&& docker-compose down
