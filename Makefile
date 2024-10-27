
postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Ds12345! -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root to_do_list

dropdb:
	docker exec -it postgres15 dropdb to_do_list

migratedown:
	migrate -path db/migration -database "postgresql://root:Ds12345!@localhost:5432/to_do_list?sslmode=disable" -verbose down

migrateup:
	migrate -path db/migration -database "postgresql://root:Ds12345!@localhost:5432/to_do_list?sslmode=disable" -verbose up
sqlc:
	sqlc generate