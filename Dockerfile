# build stage
FROM golang:1.14.4 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main.go

# final stage
FROM alpine:3.11.6

COPY --from=builder /app/main /app/

ENTRYPOINT ["/app/main"]
