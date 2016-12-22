#!/bin/sh

set -ex

source_dir=`dirname $0`

cp $source_dir/config/database-jenkins.yml $source_dir/config/database.yml

go install github.com/af83/edwig

export EDWIG_ENV=test
(cd $source_dir; $GOPATH/bin/edwig migrate up)

go test -v -cover github.com/af83/edwig/...

ruby_bin_dir=`ls -d /var/lib/gems/*/bin | tail -1`
bundle=$ruby_bin_dir/bundle

if [ -x $bundle ]; then
    cd $source_dir
    $bundle install --deployment --path $GOPATH/tmp
    $bundle exec license_finder
		mkdir -p tmp/cucumber
    $bundle exec cucumber --format json --out tmp/cucumber/cucumber.json --format html --out tmp/cucumber/cucumber.html --format pretty --no-color
else
    echo "Bundle not detected, cucumber tests are skipped"
fi
