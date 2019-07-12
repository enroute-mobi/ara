#!/bin/bash

command=${1:-api}

export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

if [ -z "$EDWIG_DB_PASSWORD" ]; then
    echo "No database password defined by EDWIG_DB_PASSWORD"
    exit 1
fi

if [ -z "$EDWIG_API_KEY" ]; then
    echo "No api key defined by EDWIG_API_KEY"
    exit 1
fi

cat > config/database.yml <<EOF
${EDWIG_ENV:-production}:
  name: ${EDWIG_DB_NAME:-edwig}
  user: ${EDWIG_DB_USER:-edwig}
  host: ${EDWIG_DB_HOST:-db}
  password: ${EDWIG_DB_PASSWORD}
  port: ${EDWIG_DB_PORT:-5432}
EOF

# echo "Current database config"
# echo "---"
# cat config/database.yml
# echo "---"

cat > config/config.yml <<EOF
syslog: ${EDWIG_SYSLOG:-false}
debug: ${EDWIG_DEBUG:-false}
apikey: ${EDWIG_API_KEY}
EOF

if [ -n "$EDWIG_LOGSTASH" ]; then
    echo "logstash: ${EDWIG_LOGSTASH}" >> config/config.yml
fi

# echo "Current Edwig config"
# echo "---"
# cat config/config.yml
# echo "---"

touch config/production.yml

echo "Start $command"
case $command in
  api)
    exec ./edwig api -listen 0.0.0.0:8080
    ;;
  shell)
    exec bash
    ;;
  migrate)
    exec ./edwig migrate up
    ;;
  *)
    exec $@
esac
