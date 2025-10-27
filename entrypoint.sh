set -e

echo "Generando clave de host SSH..."
ssh-keygen -t rsa -b 4048 -f ./ssh_host_key -N ""

echo "Iniciando servidor del juego..."
exec ./server
