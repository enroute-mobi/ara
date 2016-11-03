# Edwig

## Installation

### Application

```
git clone git@github.com:af83/edwig.git
cd edwig
```

### Postgres

```
psql
postgres=# CREATE USER "edwig" SUPERUSER PASSWORD 'edwig';
CREATE ROLE
postgres=# CREATE DATABASE edwig;
CREATE DATABASE
postgres=# CREATE DATABASE edwig_test;
CREATE DATABASE
```
Database configuration can be defined in `config/database.yml`

#### Apply migrations
```
edwig migrate up
```

#### Rollback migrations
```
edwig migrate down
```

Migration folder can be specified with
```
edwig migrate -path=path/to/migrate/files [...]
```
#### Populate
```
psql -U edwig -d edwig -a -f model/populate.sql
```

## Run Edwig

### Server
```
edwig api
```

### Checkstatus
```
edwig check http://url.to.check
```

### Configuration

To run Edwig in a specific environment
```
EDWIG_ENV=development edwig [...]
```
default environment is `development`

To load a custom configuration directory
```
edwig -config=config/directory/path [...]
```
You need to specify 3 files : `config.yml`, `database.yml`, `<environment>.yml`

Configuration try the given path if it exists, then the environment variable EDWIG_CONFIG, then `/config`