setup:
	go mod tidy

clean:
	rm -rf ./bin

build-prof: clean setup
	CGO_ENABLED=0 go build -pgo=./cmd/cpu.pprof -v -o ./bin/rinha ./cmd/rinha.go

build: clean setup
	go build -v -o ./bin/rinha ./cmd/rinha.go

image:
	docker build -t suricat89/rinha-2024-q1 .

docker-push:
	docker buildx build --push --platform linux/amd64 --tag suricat/rinha-2024-q1 .

dev-postgres:
	docker compose -f ./docker-compose.postgres.yml up db_postgres redis

dev-mongodb:
	docker compose -f ./docker-compose.mongodb.yml up db_mongo redis

image-cleanup:
	@for i in $$(docker ps --filter "name=rinha-go-suricat-api01" --filter "name=rinha-go-suricat-api02" --format "{{.ID}}"); do docker rm -f $$i; done
	@for i in $$(docker image ls --filter "reference=rinha-go-suricat-api01" --filter "reference=rinha-go-suricat-api02" --format "{{.ID}}"); do docker image rm -f $$i; done

run-postgres: image-cleanup
	docker compose -f ./docker-compose.postgres.yml up

run-mongodb: image-cleanup
	docker compose -f ./docker-compose.mongodb.yml up

stats-postgres:
	docker container stats rinha-go-suricat-db_postgres-1 rinha-go-suricat-api01-1 rinha-go-suricat-api02-1 rinha-go-suricat-nginx-1 rinha-go-suricat-redis-1

stats-mongo:
	docker container stats rinha-go-suricat-db_mongo-1 rinha-go-suricat-api01-1 rinha-go-suricat-api02-1 rinha-go-suricat-nginx-1 rinha-go-suricat-redis-1
