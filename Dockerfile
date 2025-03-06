FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN apk --no-cache add ca-certificates

RUN go build -o ./golang-server ./cmd


FROM alpine:latest AS runner

WORKDIR /app
COPY --from=builder /app/golang-server .

EXPOSE 8080
ENTRYPOINT ["./golang-server"]