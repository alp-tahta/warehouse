up:
	docker build -t warehouse:multistage -f Dockerfile.multistage .
	docker run -d -p 8080:8080 --name warehouse warehouse:multistage

down:
	docker container stop warehouse
	docker container rm warehouse
	docker image rm warehouse:multistage

rip:
	docker system prune -a

sal:
	docker compose build --no-cache warehouse
	docker compose up --no-deps warehouse

.PHONY: up down rip sal