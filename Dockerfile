# builder
FROM golang AS builder

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd

# base image

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

COPY --from=builder /go/bin/migrate /usr/local/bin/
COPY migrations ./migrations
COPY templates ./templates

ENTRYPOINT ["./app"]

#EXPOSE 8080