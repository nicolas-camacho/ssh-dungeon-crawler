#!/bin/sh
set -e

echo "=== Starting SSH Dungeon Crawler Server ==="

KEY_FILE="./ssh_host_key"

if [ ! -f "$KEY_FILE" ]; then
    echo "Host key not found. Generating a new one..."
    ssh-keygen -t rsa -b 4096 -f "$KEY_FILE" -N "" -C "ssh-dungeon-crawler"
    chmod 600 "$KEY_FILE"
    echo "Host key generated successfully."
else
    echo "Existing host key found."
    chmod 600 "$KEY_FILE"
fi

echo "Configuration:"
echo "  SSH_HOST: ${SSH_HOST:-0.0.0.0}"
echo "  SSH_PORT: ${SSH_PORT:-2222}"
echo ""

if [ ! -f "./server" ]; then
    echo "ERROR: Server binary not found!"
    exit 1
fi

if [ ! -d "./data" ]; then
    echo "WARNING: Game data directory not found!"
fi

echo "Starting game server in SSH mode..."
echo "Players can connect using: ssh -p ${SSH_PORT:-2222} ${SSH_HOST:-0.0.0.0}"
echo ""

exec ./server -ssh
