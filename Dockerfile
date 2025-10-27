
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM alpine:latest

RUN apk add --no-cache openssh-server openssh-keygen && \
    mkdir -p /var/run/sshd && \
    chmod 755 /var/run/sshd

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/data ./data
COPY entrypoint.sh .

RUN chmod +x ./entrypoint.sh && \
    chmod +x ./server

EXPOSE 2222

ENV SSH_HOST=0.0.0.0
ENV SSH_PORT=2222

CMD ["./entrypoint.sh"]
