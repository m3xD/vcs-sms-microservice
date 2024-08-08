postgres:
	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=2108 -p 5432:5432 -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root sms

dropdb:
	docker exec -it postgres dropdb sms

migrateup:
	migrate -path db/migration -database "postgresql://root:2108@localhost:5432/sms?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:2108@localhost:5432/sms?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go vcssms/db/sqlc Store      
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock