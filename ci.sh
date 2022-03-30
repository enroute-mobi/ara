#!/bin/sh

set -ex

source_dir="$(pwd -P)/$(dirname "$0")"

cat > "$source_dir/config/database.yml" <<EOF
test:
  name: ${ARA_DB_NAME:-ara_test}
  user: ${ARA_DB_USER:-jenkins}
  host: ${ARA_DB_HOST:-db}
  password: ${ARA_DB_PASSWORD}
  port: ${ARA_DB_PORT:-5432}
EOF

# go install honnef.co/go/tools/cmd/staticcheck@latest

cd $source_dir

# staticcheck ./...

go install -v ./...

export ARA_ENV=test
export ARA_ROOT=$source_dir
"$GOPATH/bin/ara" migrate up

go install github.com/schrej/godacov@latest

go test -coverprofile=coverage.out -p 1 ./...

if [ -n "$CODACY_PROJECT_TOKEN" ]; then
    $GOPATH/bin/godacov -t "$CODACY_PROJECT_TOKEN" -r ./coverage.out -c "$BITBUCKET_COMMIT"
fi

cd "$source_dir"
# $bundle exec license_finder
bundle exec bundle-audit check --update
bundle exec cucumber --tags "not @wip" --publish
