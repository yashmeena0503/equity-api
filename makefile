.PHONY: api sqlc sim web db dbx x

all: dev

# TODO: add api-dev maybe
dev:
	docker compose up -d mariadb adminer api webapp-dev

# api
api:
	cd api && go build -o tmp/ && ./tmp/api
# sqlc generate, .sql -> .go
sqlc:
	cd api && nix develop -c go tool sqlc generate

# sim
sim:
	cd sim && nix develop -c go build -o tmp/ && ./tmp/sim

# webapp
web:
	cd webapp && nix develop -c npm run dev

# misc
db:
	docker compose up -d adminer
dbx:
	docker compose down mariadb adminer
x:
	docker compose down
