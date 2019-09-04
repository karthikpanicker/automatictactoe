#!/bin/bash
# The script stops running containers , deletes earlier image, pull and runs latest image

DOCKER_IMAGE=$1
BUILD_DIR=$2
ETSY_CONSUMER_KEY=$3
ETSY_CONSUMER_SECRET=$4
TRELLO_CONSUMER_KEY=$5
TRELLO_CONSUMER_SECRET=$6
executor()
{
cd $BUILD_DIR
pwd
docker-compose down
docker rmi -f $(docker images | grep $DOCKER_IMAGE | tr -s ' ' | cut -d ' ' -f 3)
docker images
docker-compose pull --quiet
ETSY_CONSUMER_KEY=$ETSY_CONSUMER_KEY ETSY_CONSUMER_SECRET=$ETSY_CONSUMER_SECRET TRELLO_CONSUMER_KEY=$TRELLO_CONSUMER_KEY TRELLO_CONSUMER_SECRET=$TRELLO_CONSUMER_SECRET docker-compose up -d --force-recreate
}

executor

