.PHONY: build run migrate clean

build:
	docker compose build

run:
	docker compose up

migrate-up:
	docker compose run app go run migrate/migrate.go 

clean:
	docker compose down -v
	