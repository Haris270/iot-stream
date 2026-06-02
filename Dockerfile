FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o engine .

# ----------- The Runner ----------- #

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/engine .

CMD ["./engine"]
