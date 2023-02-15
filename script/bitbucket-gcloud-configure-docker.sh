#!/bin/bash -e

# Manage base64 encoded or raw GCLOUD_API_KEYFILE
if echo "$GCLOUD_API_KEYFILE" | grep -q "private_key"; then
    echo -En "$GCLOUD_API_KEYFILE" > ~/.gcloud-api-key.json
else
    echo "$GCLOUD_API_KEYFILE" | base64 -d > ~/.gcloud-api-key.json
fi

gcloud auth activate-service-account --key-file ~/.gcloud-api-key.json
gcloud config set project "$GCLOUD_PROJECT"
gcloud auth configure-docker --quiet
