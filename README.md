# SSH Dungeon Crawler

¡Un simple juego de exploración de mazmorras (Dungeon Crawler) al que juegas a través de SSH!

Este proyecto es una aplicación de Go que utiliza la suite de herramientas de Charmbracelet para crear un servidor SSH que aloja un juego interactivo de mazmorras basado en texto. Los jugadores pueden conectarse al servidor y navegar por un mapa generado proceduralmente, encontrar diferentes tipos de salas y ver sus estadísticas.

## Tecnologías Utilizadas

-   **Go**: El lenguaje de programación principal.
-   **Charmbracelet Wish**: Para crear el servidor SSH.
-   **Charmbracelet Bubble Tea**: Para construir la interfaz de usuario interactiva (TUI).
-   **Charmbracelet Lip Gloss**: Para estilizar la salida en la terminal.

## Cómo Empezar

### Prerrequisitos

-   Tener Go (versión 1.25 o superior) instalado.
-   Un cliente SSH.

### Ejecutar el Servidor

1.  **Clona el repositorio:**
    ```bash
    git clone <URL_DEL_REPOSITORIO>
    cd ssh-dungeon-crawler
    ```

2.  **Instala las dependencias:**
    ```bash
    go mod tidy
    ```

3.  **Ejecuta el servidor:**
    ```bash
    go run .
    ```
    El servidor SSH comenzará a escuchar en `localhost:2222`.

### Conectarse al Juego

1.  Abre una nueva terminal y conéctate al servidor usando SSH:
    ```bash
    ssh localhost -p 2222
    ```

2.  ¡El juego comenzará automáticamente!

## Controles

-   **Movimiento**: Usa las **teclas de flecha** o las teclas **W, A, S, D** para mover a tu personaje por el mapa.
-   **Salir**: Presiona **q** o **Ctrl+C** para salir del juego.

## Estructura del Proyecto

```
.
├── .gitignore
├── go.mod
├── go.sum
├── main.go         # Punto de entrada, configuración y ejecución del servidor SSH
├── README.md
└── game/
    ├── app.go      # Aplicación principal de Bubble Tea y gestión de estados
    ├── gameplay.go # Lógica del juego principal, movimiento y renderizado
    ├── loading.go  # Lógica y renderizado de la pantalla de carga
    ├── map.go      # Generación procedural del mapa
    ├── menu.go     # Lógica y renderizado del menú principal
    └── model.go    # Estructuras de datos y modelos del juego
```