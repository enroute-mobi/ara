#!/bin/bash

command=${1:-api}

export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

if [ -z "$ARA_DB_PASSWORD" ]; then
    echo "No database password defined by ARA_DB_PASSWORD"
    exit 1
fi

if [ -z "$ARA_API_KEY" ]; then
    echo "No api key defined by ARA_API_KEY"
    exit 1
fi

if [ -n "$GCLOUD_KEYFILE_JSON" ]; then
    export GOOGLE_APPLICATION_CREDENTIALS="gcloud-keyfile.json"
    echo -En "$GCLOUD_KEYFILE_JSON" > "$GOOGLE_APPLICATION_CREDENTIALS"
    unset GCLOUD_KEYFILE_JSON
fi

cat > config/database.yml <<EOF
${ARA_ENV:-production}:
  name: ${ARA_DB_NAME:-ara}
  user: ${ARA_DB_USER:-ara}
  host: ${ARA_DB_HOST:-db}
  password: ${ARA_DB_PASSWORD}
  port: ${ARA_DB_PORT:-5432}
EOF

# echo "Current database config"
# echo "---"
# cat config/database.yml
# echo "---"

cat > config/config.yml <<EOF
syslog: ${ARA_SYSLOG:-false}
debug: ${ARA_DEBUG:-false}
apikey: ${ARA_API_KEY}
bigqueryprojectid: "${GCLOUD_PROJECT}"
bigquerydataset: "${ARA_BIGQUERY_DATASET}"
bigquerytable: "exchange_events"
EOF

if [ -n "$ARA_LOGSTASH" ]; then
    echo "logstash: ${ARA_LOGSTASH}" >> config/config.yml
fi

# echo "Current Ara config"
# echo "---"
# cat config/config.yml
# echo "---"

touch config/production.yml

echo "Start $command"
case $command in
  api)
    if [ "$RUN_MIGRATIONS" = "true" ]; then
        ./ara migrate up || exit $?
    fi
    exec ./ara api -listen 0.0.0.0:8080
    ;;
  shell)
    exec bash
    ;;
  migrate)
    exec ./ara migrate up
    ;;
  *)
    exec $@
esac
