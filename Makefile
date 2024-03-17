migrationup:
	migrate -path db/migration -database "postgresql://postgres:147563@localhost:5432/tech_school_course?sslmode=disable" -verbose up

migrationdown:
	migrate -path db/migration -database "postgresql://postgres:147563@localhost:5432/tech_school_course?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

.PHONY: migrationup migrationdown sqlc test server mockgen