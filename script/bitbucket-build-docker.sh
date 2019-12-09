#!/bin/sh -e

export IMAGE_NAME="eu.gcr.io/$GCLOUD_PROJECT/$BITBUCKET_REPO_SLUG:$BITBUCKET_COMMIT"

# Build image
docker build . -t "$IMAGE_NAME" --build-arg VERSION="build-$BITBUCKET_BUILD_NUMBER"

# Publish image
docker push "$IMAGE_NAME"
