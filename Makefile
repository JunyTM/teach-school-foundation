migrationup:
	migrate -path db/migration -database "postgresql://postgres:147563@localhost:5432/tech_school_course?sslmode=disable" -verbose up

migrationdown:
	migrate -path db/migration -database "postgresql://postgres:147563@localhost:5432/tech_school_course?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: migrationup migrationdown sqlc