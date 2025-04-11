up:
	docker compose up -d

down:
	docker compose down -v

rip:
	docker system prune -a

sal:
	docker compose build --no-cache warehouse
	docker compose up --no-deps warehouse

.PHONY: up down rip sal