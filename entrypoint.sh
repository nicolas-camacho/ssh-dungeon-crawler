#!/bin/sh
# entrypoint.sh
set -e

KEY_FILE="./ssh_host_key"

if [ ! -f "$KEY_FILE" ]; then
    echo "Host key not found. Generating a new one..."
    ssh-keygen -t rsa -b 4096 -f "$KEY_FILE" -N ""
    echo "Host key generated successfully."
else
    echo "Existing host key found."
fi

echo "Starting game server..."
exec ./server
