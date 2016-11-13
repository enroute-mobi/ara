TMPGOPATH := build
GOPKG := github.com/af83/edwig

build:
	mkdir -p "${TMPGOPATH}/src/${GOPKG}"
	find . -path ./debian -prune -o -type f -exec tar cf - {} + | (cd "${TMPGOPATH}/src/${GOPKG}" && tar xf -)
	GOPATH=${PWD}/${TMPGOPATH} go install ${GOPKG}

install:
	cp ${TMPGOPATH}/bin/edwig ${DESTDIR}/usr/bin

clean:
	rm -rf build
	rm Makefile
