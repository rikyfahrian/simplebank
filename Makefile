postgresrun:
	docker run --name techschoolpgsql -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it techschoolpgsql createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it techschoolpgsql dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
	ke
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc: 
	sqlc generate
test:
	go test -v ./...
run: 
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go -source db/sqlc/store.go 
	