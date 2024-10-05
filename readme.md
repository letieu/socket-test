# Socket Test

This project is a simple WebSocket server implemented in Go. It serves HTML files and handles WebSocket connections for clients and monitors.

## Features

- Serve HTML files (`index.html` and `client.html`) using Go's `embed` package.
- Handle WebSocket connections for clients and monitors.
- Broadcast messages to connected clients and monitors.
- Support for connecting and disconnecting events.

## Prerequisites

- Go 1.16 or later
- [Gorilla WebSocket](https://github.com/gorilla/websocket) package

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/socket-test.git
    cd socket-test
    ```

## Usage

1. Set the environment variable for the server port (optional):
    ```sh
    export S_PORT=8082
    ```

2. Run the server:
    ```sh
    go run main.go
    ```

3. Open your browser and navigate to `http://localhost:8082` to access the home page.

## Endpoints

- `/` - Serves the `index.html` file.
- `/client` - Serves the `client.html` file.
- `/ws` - Handles WebSocket connections.

## WebSocket Events

- `connect` - Triggered when a client or monitor connects.
- `disconnect` - Triggered when a client or monitor disconnects.
- `message` - Handles messages sent by clients or monitors.

## Running Tests

To run tests, use the following command:
```sh
go test ./...
```
