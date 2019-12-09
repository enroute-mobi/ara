#!/bin/sh -e

TAG=$1

# Tag image in registry with given label
IMAGE_NAME="eu.gcr.io/$GCLOUD_PROJECT/$BITBUCKET_REPO_SLUG:$BITBUCKET_COMMIT"
TAGGED_IMAGE_NAME="eu.gcr.io/$GCLOUD_PROJECT/$BITBUCKET_REPO_SLUG:$TAG"

echo "Tag $IMAGE_NAME as $TAGGED_IMAGE_NAME"

docker pull "$IMAGE_NAME"
docker tag "$IMAGE_NAME" "$TAGGED_IMAGE_NAME"
docker push "$TAGGED_IMAGE_NAME"
