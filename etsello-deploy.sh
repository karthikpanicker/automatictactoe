#!/bin/bash
# The script stops running containers , deletes earlier image, pull and runs latest image

DOCKER_IMAGE=$1
BUILD_DIR=$2
TRELLO_CONSUMER_SECRET=$3
executor()
{
cd $BUILD_DIR
pwd
docker-compose down
docker rmi -f $(docker images | grep $DOCKER_IMAGE | tr -s ' ' | cut -d ' ' -f 3)
docker images
docker-compose pull --quiet
TRELLO_CONSUMER_SECRET=$TRELLO_CONSUMER_SECRET docker-compose up -d --force-recreate
}

executor

