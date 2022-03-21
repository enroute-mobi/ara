export ARA_ROOT=$(PWD)
export ARA_CONFIG=$(PWD)/config

run: # for dev
	go run -race ara.go api

migrations:
	go run ara.go migrate up
	ARA_ENV=test go run ara.go migrate up

rollback_migrations:
	go run ara.go migrate down
	ARA_ENV=test go run ara.go migrate down

populate:
	psql -U ara -d ara -a -f model/populate.sql

tests:
	go test bitbucket.org/enroute-mobi/ara/... -p 1 -count 1

cucumber:
	bundle exec cucumber -t 'not @wip'

gen_gtfsrt_bindings:
	wget https://developers.google.com/transit/gtfs-realtime/gtfs-realtime.proto
	protoc --go_out=. --go_opt=Mgtfs-realtime.proto=gtfs/ gtfs-realtime.proto
	rm gtfs-realtime.proto