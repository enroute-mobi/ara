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

cd $source_dir

go install -v ./...

export EDWIG_ENV=test
$GOPATH/bin/edwig migrate up

go test -p 1 ./...

tmp_dir=$GOPATH/tmp

cd $source_dir
# $bundle exec license_finder
bundle exec bundle-audit check --update

mkdir -p $tmp_dir/cucumber
bundle exec cucumber --tags "~@wip" --format json --out $tmp_dir/cucumber/cucumber.json --format html --out $tmp_dir/cucumber/index.html --format pretty --no-color
