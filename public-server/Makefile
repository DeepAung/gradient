air:
	air -c .air.toml
tailwind:
	asfd

db.create:
	docker create --name gradient-postgresql --env-file postgres.env.dev -p 5432:5432 bitnami/postgresql:latest
db.start:
	docker start -i gradient-postgresql

db.migrate:
	docker cp ./migrations/migrate.sql gradient-postgresql:/migrate.sql
	docker exec -it gradient-postgresql psql postgresql://myuser:mypassword@0.0.0.0/mydb -f /migrate.sql
db.seed:
	docker cp ./migrations/seed.sql gradient-postgresql:/seed.sql
	docker exec -it gradient-postgresql psql postgresql://myuser:mypassword@0.0.0.0/mydb -f /seed.sql
	go run migrations/seed.go
db.reset:
	make db.migrate db.seed
