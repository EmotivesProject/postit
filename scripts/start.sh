#!/bin/bash
# Meant to be ran from root directory

stage="$1"
build="$2"
export GMAIL_PASSWORD="$3"
export NOTIFICATION_AUTH="$4"

if [ "$stage" == 'prod' ]; then
	if [ "$build" == 'yes' ]; then
		echo Building prod and running it
		docker-compose -f docker/postit/docker-compose.prod.yml up -d --build
		exit
	fi
	echo Running prod
	docker-compose -f docker/postit/docker-compose.prod.yml up -d
	exit
fi

if [ "$build" == 'yes' ]; then
	echo Building dev and running it
	docker-compose -f docker/postit/docker-compose.yml up -d --build
	exit
fi

echo Running dev
docker-compose -f docker/postit/docker-compose.yml up -d

