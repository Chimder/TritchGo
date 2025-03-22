FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN apk --no-cache add ca-certificates

RUN go build -o main ./cmd


FROM alpine:latest AS runner

WORKDIR /app
COPY --from=builder /app/maim .

EXPOSE 8080
ENTRYPOINT ["/main"]

# FROM golang:1.24.1 AS builder
# WORKDIR /app


# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# COPY . .
# RUN CGO_ENABLED=0 go build -o main ./cmd

# FROM gcr.io/distroless/static-debian12

# COPY --from=builder /app/main /

# EXPOSE 8080
# CMD ["/main"]

