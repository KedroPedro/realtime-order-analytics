FROM docker.io/library/golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/order-service/

FROM alpine:3.20

WORKDIR /root/

COPY --from=builder /app/main .

ENV HTTP_ADDR=0.0.0.0:8080
ENV POSTGRES_CONN_STRING=postgresql://user:password@postgresql:5432/orderdb

EXPOSE 8080

CMD ["./main"]