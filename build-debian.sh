#!/bin/bash

set -e

TARGET_DIR=${1:-target}

if [ -z "$BUILD_BRANCH" -o "$BUILD_BRANCH" == "master" ]; then
  deploy_env="dev"
else
  deploy_env=$BUILD_BRANCH
fi

if [ -n "$BUILD_NUMBER" ]; then
    debian_release=`dpkg-parsechangelog --count 1 | sed -n '/^Version: / s/Version: //p'`
    dch --changelog debian/changelog --newversion "${debian_release}+${deploy_env}${BUILD_NUMBER}" --distribution unstable 'New build'
fi

mk-build-deps -i -t 'apt-get -o Debug::pkgProblemResolver=yes --no-install-recommends -y'
debuild -us -uc

mkdir -p $TARGET_DIR
cp ../*.deb $TARGET_DIR
