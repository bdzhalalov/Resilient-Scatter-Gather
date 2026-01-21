FROM golang:1.25-alpine AS builder

RUN apk --no-cache add bash make gcc musl-dev

WORKDIR /var/www/

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o ./cmd main.go

FROM alpine AS runner

COPY --from=builder /var/www/cmd /
COPY config /config

WORKDIR /var/www/

EXPOSE 8080

CMD ["/cmd"]
