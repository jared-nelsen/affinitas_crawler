
run:
	go build
	./cabo_affinitas

database-up:
	docker compose up

database-up-rebuild:
	docker compose up --build

database-down:
	docker compose down

kill-db-port:
	sudo fuser -k 5432/tcp