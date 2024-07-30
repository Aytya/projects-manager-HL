.PHONY: build run

build:
	docker-compose build

down:
	docker-compose down

test:
	go test github.com/Aytya/projects-manager-HL/tests

run:
	docker-compose up