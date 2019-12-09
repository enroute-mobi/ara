#!/bin/sh -e

gcloud container clusters get-credentials "$GCLOUD_CLUSTER" --zone "$GCLOUD_ZONE"

IMAGE_NAME="eu.gcr.io/$GCLOUD_PROJECT/$BITBUCKET_REPO_SLUG:$BITBUCKET_COMMIT"
kubectl set image deployment --namespace="$GCLOUD_NAMESPACE" api api="$IMAGE_NAME" --record
