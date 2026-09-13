FROM golang:1.27-alpine AS builder

WORKDIR /analytic

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o main ./cmd/analytic/

FROM alpine:3.20

WORKDIR /root/

COPY --from=builder /analytic/main .

EXPOSE 8080

CMD ["./main"]

