FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/server .

FROM alpine:latest

RUN apk add --no-cache openssh

WORKDIR /app

COPY --from=builder /app/server .
COPY ./data ./data
COPY entrypoint.sh .

RUN chmod +x ./entrypoint.sh

CMD ["./entrypoint.sh"]
