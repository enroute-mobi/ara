PWD := $(shell pwd)
TMPGOPATH := "${PWD}/go"
GOPKG := github.com/af83/edwig
VERSION := $(shell head -1 debian/changelog | sed 's/edwig (\(.*\)) .*/\1/')
LDFLAGS := -ldflags "-X github.com/af83/edwig/version.value=${VERSION}"

export PATH := /go/bin:$(PATH)

build:
	GOROOT=/go GOPATH=${TMPGOPATH} GO111MODULE=on go install ${LDFLAGS} ${GOPKG}

install:
	cp ${TMPGOPATH}/bin/edwig ${DESTDIR}/usr/bin

clean:
	rm -rf ${TMPGOPATH}
	rm Makefile
