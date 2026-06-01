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

redis_tests:
	ARA_REDIS_ADDR=127.0.0.1:6379 go test -coverprofile=coverage.out -p 1 -count 1  ./...

cucumber:
	go build && bundle exec cucumber -t 'not @wip'

redis_cucumber:
	go build && ARA_REDIS_ADDR=127.0.0.1:6379 bundle exec cucumber -t 'not @wip'

gen_proto_bindings:
	protoc --go_out=. --go_opt=paths=source_relative \
		audit/exchangepb/exchange.proto \
		audit/partnerpb/partner.proto \
		audit/vehiclepb/vehicle.proto \
		audit/controlpb/control.proto \
		audit/longtermsvpb/longtermsv.proto

gen_gtfsrt_bindings:
	wget https://raw.githubusercontent.com/google/transit/refs/heads/master/gtfs-realtime/proto/gtfs-realtime.proto
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
	install -t build/config -D config/config.yml config/database.yml config/test.yml
	mkdir -p build/config

clean:
	rm -rf build
