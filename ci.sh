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
    $bundle install --deployment
    $bundle exec cucumber
else
    echo "Bundle not detected, cucumber tests are skipped"
fi
