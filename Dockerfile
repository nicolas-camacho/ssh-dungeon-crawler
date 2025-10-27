FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copiar módulos y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Imagen final
FROM alpine:latest

# Instalar openssh, bash y dependencias para mejor soporte de terminal
RUN apk add --no-cache openssh-server openssh-keygen bash ncurses-terminfo && \
    mkdir -p /var/run/sshd && \
    chmod 755 /var/run/sshd

WORKDIR /app

# Copiar binario y archivos necesarios
COPY --from=builder /app/server .
COPY --from=builder /app/data ./data

# Dar permisos de ejecución
RUN chmod +x ./server

# Exponer puerto SSH
EXPOSE 2222

# Variables de entorno por defecto
ENV SSH_HOST=0.0.0.0
ENV SSH_PORT=2222

# Ejecutar directamente con bash inline
CMD ["/bin/sh", "-c", "\
    echo '=== Starting SSH Dungeon Crawler Server ===' && \
    KEY_FILE='./ssh_host_key' && \
    if [ ! -f \"$KEY_FILE\" ]; then \
    echo 'Generating SSH host key...' && \
    ssh-keygen -t rsa -b 4096 -f \"$KEY_FILE\" -N '' -C 'ssh-dungeon-crawler' && \
    chmod 600 \"$KEY_FILE\" && \
    echo 'Host key generated successfully.'; \
    else \
    echo 'Existing host key found.' && \
    chmod 600 \"$KEY_FILE\"; \
    fi && \
    echo \"SSH_HOST: $SSH_HOST\" && \
    echo \"SSH_PORT: $SSH_PORT\" && \
    echo 'Starting game server in SSH mode...' && \
    exec ./server -ssh \
    "]
