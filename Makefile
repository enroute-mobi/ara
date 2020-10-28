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