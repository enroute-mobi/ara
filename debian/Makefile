TMPGOPATH := build
GOPKG := github.com/af83/edwig
CMD := head -1 debian/changelog | sed 's/edwig (\(.*\)) .*/\1/'
VERSION := $(shell head -1 debian/changelog | sed 's/edwig (\(.*\)) .*/\1/')
LDFLAGS := -ldflags "-X github.com/af83/edwig/version.value=${VERSION}"

build:
	mkdir -p "${TMPGOPATH}/src/${GOPKG}"
	find . -path ./debian -prune -o -path ./.git -prune -o -type f -exec tar cf - {} + | (cd "${TMPGOPATH}/src/${GOPKG}" && tar xf -)
	GOPATH=${PWD}/${TMPGOPATH} go install ${LDFLAGS} ${GOPKG}

install:
	cp ${TMPGOPATH}/bin/edwig ${DESTDIR}/usr/bin

clean:
	rm -rf build
	rm Makefile
