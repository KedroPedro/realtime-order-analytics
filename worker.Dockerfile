FROM golang:1.27-alpine AS builder

WORKDIR /worker

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o main ./cmd/worker/

FROM alpine:3.20

WORKDIR /root/

COPY --from=builder /worker/main .

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]
