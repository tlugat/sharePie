FROM golang:1.18 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shapePie-api .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/shapePie-api .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./shapePie-api"]
