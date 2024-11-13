export ARA_ROOT=$(PWD)
export ARA_CONFIG=$(PWD)/config

run: # for dev
	go run -race ara.go api

convert:
	go run ara.go convert $(SCHEMA_NAME)

test_migrations:
	ARA_ENV=test go run ara.go migrate up

dev_migrations:
	go run ara.go migrate up

migrations: dev_migrations test_migrations

rollback_migrations:
	go run ara.go migrate down
	ARA_ENV=test go run ara.go migrate down

populate:
	psql -U ara -d ara -a -f model/populate.sql

tests:
	go test -coverprofile=coverage.out -p 1 -count 1  ./...

cucumber:
	bundle exec cucumber -t 'not @wip'

gen_gtfsrt_bindings:
	wget https://raw.githubusercontent.com/google/transit/master/gtfs-realtime/proto/gtfs-realtime.proto
	protoc --go_out=. --go_opt=Mgtfs-realtime.proto=gtfs/ gtfs-realtime.proto
	rm gtfs-realtime.proto

go_dependencies:
	go mod vendor

ruby_dependencies:
	MAKE="make --jobs $(nproc)" bundle install --jobs `nproc` --path vendor/bundle

dependencies: go_dependencies ruby_dependencies

build:
	go install -v ./...
	mkdir -p build
	install --mode=+x ${GOPATH}/bin/ara build
	install -t build/db/migrations -D db/migrations/*.sql
	install -t build/siri/templates -D siri/templates/*.template
	mkdir -p build/config

clean:
	rm -rf build
