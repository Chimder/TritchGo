generate:
	cd proto && protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./**/*.proto

gpweb:
	pgweb --url "postgresql://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"

protoui:
	grpcui -plaintext localhost:50051
