include .env
export $(shell sed 's/=.*//' .env)

up:
	goose -dir db/migration postgres ${DB_URL} up

create_%:
	goose -dir db/migration create $* sql

generate:
	cd proto && protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./**/*.proto

gpweb:
	pgweb --url ${DB_URL}

protoui:
	grpcui -plaintext localhost:50051

apply:
	kubectl apply -f k8s/app.yaml
	kubectl delete secret app-secrets --ignore-not-found
	kubectl create secret generic app-secrets --from-env-file=.env

restart:
	kubectl rollout restart deployment golang-server

getPostgresIp:
	docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' my-postgres

ipwsl:
	ip addr show eth0 | grep 'inet ' | awk '{print $2}' | cut -d'/' -f1

wrk:
	wrk -t2 -c100 -d30s http://localhost:8080/go-json-gzip