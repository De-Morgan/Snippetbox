
composeup:
	docker compose up

composedown:
	docker compose down
	
app:
	go run ./cmd/web

all:
	make composeup && make app