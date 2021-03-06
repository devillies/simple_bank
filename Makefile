dockerdb:
	docker stop adminPanel && docker start postgres12

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=slot123 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:slot123@localhost:5432/simple_bank?sslmode=disable" --verbose up 

migrateup1:
	migrate -path db/migration -database "postgresql://root:slot123@localhost:5432/simple_bank?sslmode=disable" --verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:slot123@localhost:5432/simple_bank?sslmode=disable" --verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:slot123@localhost:5432/simple_bank?sslmode=disable" --verbose down 1

sqlc:
	sqlc generate


test:
	go test -v -cover ./...


server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/devillies/simple_bank/db/sqlc Store 

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc server mock dockerdb