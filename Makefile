migrations:
	go run edwig.go migrate up
	EDWIG_ENV=test go run edwig.go migrate up

rollback_migrations:
	go run edwig.go migrate down
	EDWIG_ENV=test go run edwig.go migrate down

populate:
	psql -U edwig -d edwig -a -f model/populate.sql

run:
	go run edwig.go api

tests:
	go test github.com/af83/edwig/... -p 1 -count 1

cucumber:
	bundle exec cucumber -t ~@wip