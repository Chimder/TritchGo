generate:
	cd proto && protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./**/*.proto

gpweb:
	pgweb --url "postgresql://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"

protoui:
	grpcui -plaintext localhost:50051

ipwsl:
	ip addr show eth0 | grep 'inet ' | awk '{print $2}' | cut -d'/' -f1

wrk:
	wrk -t2 -c100 -d30s http://localhost:8080/go-json-gzip