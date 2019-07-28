#!/bin/sh

set -e

TARGET_DIR=${1:-target}

mk-build-deps -i -t 'apt-get -o Debug::pkgProblemResolver=yes --no-install-recommends -y'
debuild -us -uc

mkdir -p $TARGET_DIR
cp ../*.deb $TARGET_DIR
