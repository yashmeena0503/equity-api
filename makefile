dev:
	cd api && go build -o tmp && ./tmp/api

# sqlc generate, .sql -> .go
generate:
	go tool sqlc generate

db:
	docker compose up -d adminer
dbx:
	docker compose down mariadb adminer

x:
	docker compose down
