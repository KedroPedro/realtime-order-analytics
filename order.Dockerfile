FROM golang:1.27-alpine AS builder

WORKDIR /order

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o main ./cmd/order-service/

FROM alpine:3.20

WORKDIR /root/

COPY --from=builder /order/main .

EXPOSE 8080

CMD ["./main"]
