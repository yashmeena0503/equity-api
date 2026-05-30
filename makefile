dev:
	cd api && go build -o tmp && ./tmp/api

db:
	docker compose up -d adminer
dbx:
	docker compose down adminer

x:
	docker compose down
