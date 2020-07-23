PWD := $(shell pwd)
TMPGOPATH := "${PWD}/go"
GOPKG := bitbucket.org/enroute-mobi/ara
VERSION := $(shell head -1 debian/changelog | sed 's/ara (\(.*\)) .*/\1/')
LDFLAGS := -ldflags "-X bitbucket.org/enroute-mobi/ara/version.value=${VERSION}"

export PATH := /go/bin:$(PATH)

build:
	GOROOT=/go GOPATH=${TMPGOPATH} GO111MODULE=on go install ${LDFLAGS} ${GOPKG}

install:
	cp ${TMPGOPATH}/bin/ara ${DESTDIR}/usr/bin

clean:
	rm -rf ${TMPGOPATH}
	rm Makefile
