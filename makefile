# backend/Makefile
dev:
#	docker-compose up -d db redis
	docker-compose up -d db
	go run main.go

dev-down:
	docker-compose down

test:
#	docker-compose up -d db redis
	docker-compose up -d db
	go test ./...
	docker-compose down