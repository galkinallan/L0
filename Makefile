all: postgres createdb migrateup nats

postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root orders


dropdb:
	docker exec -it postgres16 dropdb --username=root orders

migrateup:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/orders?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/orders?sslmode=disable" -verbose down

nats:
	docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223


.PHONY: postgres createdb dropdb migrateup migratedown nats all