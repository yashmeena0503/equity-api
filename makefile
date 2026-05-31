dev:
	cd api && go build -o tmp && ./tmp/api

# sqlc generate, .sql -> .go
sqlc:
	cd api && nix develop -c go tool sqlc generate

db:
	docker compose up -d adminer
dbx:
	docker compose down mariadb adminer

x:
	docker compose down
