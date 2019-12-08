#!/bin/sh

set -ex

source_dir=$(dirname "$0")

cat > "$source_dir/config/database.yml" <<EOF
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
"$GOPATH/bin/edwig" migrate up

go get github.com/schrej/godacov

go test -coverprofile=coverage.out -p 1 ./...

if [ -n "$CODACY_PROJECT_TOKEN" ]; then
    $GOPATH/bin/godacov -t "$CODACY_PROJECT_TOKEN" -r ./coverage.out -c "$BITBUCKET_COMMIT"
fi

cd "$source_dir"
# $bundle exec license_finder
bundle exec bundle-audit check --update

bundle exec cucumber --tags "~@wip" --format pretty --no-color
