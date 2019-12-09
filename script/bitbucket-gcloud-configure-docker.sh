#!/bin/sh -e

echo "$GCLOUD_API_KEYFILE" | base64 -d > ~/.gcloud-api-key.json
gcloud auth activate-service-account --key-file ~/.gcloud-api-key.json
gcloud config set project "$GCLOUD_PROJECT"
gcloud auth configure-docker --quiet
