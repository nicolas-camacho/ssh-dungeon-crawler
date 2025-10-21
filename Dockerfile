FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY ./ssh_host_key ./ssh_host_key

COPY ./data ./data

CMD [ "./server" ]
