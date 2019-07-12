#!/bin/sh

set -ex

source_dir=`dirname $0`

cat > $source_dir/config/database.yml <<EOF
test:
  name: ${EDWIG_DB_NAME:-edwig_test}
  user: ${EDWIG_DB_USER:-jenkins}
  host: ${EDWIG_DB_HOST:-db}
  password: ${EDWIG_DB_PASSWORD}
  port: ${EDWIG_DB_PORT:-5432}
EOF

#(cd $source_dir; GO111MODULE=on go get -d -v ./...)
go install github.com/af83/edwig

PGHOST=${EDWIG_DB_HOST:-localhost} PGUSER=${EDWIG_DB_USER:-jenkins} PGPASSWORD=${EDWIG_DB_PASSWORD} createdb $EDWIG_DB_NAME

export EDWIG_ENV=test
(cd $source_dir; $GOPATH/bin/edwig migrate up)

go test -v -cover github.com/af83/edwig/...

tmp_dir=$GOPATH/tmp

cd $source_dir
# $bundle exec license_finder
bundle exec bundle-audit check --update

mkdir -p $tmp_dir/cucumber
bundle exec cucumber --tags "~@wip" --format json --out $tmp_dir/cucumber/cucumber.json --format html --out $tmp_dir/cucumber/index.html --format pretty --no-color
