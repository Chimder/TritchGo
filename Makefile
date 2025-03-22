include .env
export $(shell sed 's/=.*//' .env)

MGDIR = sql/migration
up:
	goose -dir $(MGDIR) postgres ${DB_URL} up

create_%:
	goose -dir $(MGDIR) create $* sql

status:
	goose -dir $(MGDIR) postgres ${DB_URL} status

reset:
	goose -dir $(MGDIR) postgres ${DB_URL} reset

down:
	goose -dir $(MGDIR) postgres ${DB_URL} down

# new-migration:
# ifndef NAME
# 	$(error Usage: make new-migration NAME=your_migration_name)
# endif
# 	goose -dir $(MGDIR) create $(NAME) sql

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

ipwsl:
	ip addr show eth0 | grep 'inet ' | awk '{print $2}' | cut -d'/' -f1

getPostgresIp:
	docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' my-postgres

# tests
vegeta:
	vegeta attack -duration=1m -rate=500 -targets=tess/vegeta/target.list -output=tess/vegeta/attack.bin

vegeta-res:
	vegeta plot -title=Attack%20Results tess/vegeta/attack.bin > tess/vegeta/results.html


wrk:
	wrk -t2 -c100 -d30s http://localhost:8080/go-json-gzip
#

stop:
	docker compose stop
dcup:
	docker compose up -d