SLEEP_TIME=0

.PHONY: lint test integration start_test_containers integration_dry get_latest_containers destory_test_containers stop_test_containers

lint:
	golangci-lint run

test:
	go test -v ./...

integration:
	make start_test_containers
	sleep $(SLEEP_TIME)
	go test -v -tags=integration ./...
	docker-compose -f docker/postit/docker-compose.test.yml down

start_test_containers:
	make get_latest_containers
	docker-compose -f docker/postit/docker-compose.test.yml down --remove-orphans
	docker-compose -f docker/postit/docker-compose.test.yml up -d --remove-orphans

stop_test_containers:
	docker-compose -f docker/postit/docker-compose.test.yml down --remove-orphans

integration_dry:
	go test -v -tags=integration ./...

get_latest_containers:
	docker pull ghcr.io/emotivesproject/postgres_db:latest
	docker pull ghcr.io/emotivesproject/uacl_api:latest

destory_test_containers:
	docker-compose -f docker/postit/docker-compose.test.yml down --remove-orphans