package main

import (
	"embed"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//go:embed index.html client.html
var embedFiles embed.FS

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]string)
var monitors = make(map[*websocket.Conn]bool)
var broadcast = make(chan Event)
var mutex = &sync.Mutex{}

type EventType string

type Event struct {
	Type    EventType `json:"type"`
	Content string    `json:"content"`
	Client  string    `json:"client,omitempty"`
	Target  string    `json:"target,omitempty"`
}

const (
	EventTypeConnect    EventType = "connect"
	EventTypeDisconnect EventType = "disconnect"
	EventTypeMessage    EventType = "message"
)

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/client", serveClient)
	http.HandleFunc("/ws", handleSocketConnections)

	go handleMessages()

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	file, err := embedFiles.ReadFile("index.html")

	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(file)
}

func serveClient(w http.ResponseWriter, r *http.Request) {
	file, err := embedFiles.ReadFile("client.html")

	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(file)
}

func handleSocketConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clientID := ws.RemoteAddr().String()
	isMonitor := r.URL.Query().Get("monitor") == "true"

	mutex.Lock()
	if isMonitor {
		monitors[ws] = true
		broadcast <- Event{Type: EventTypeConnect, Content: "Monitor connected", Client: clientID}
	} else {
		clients[ws] = clientID
		broadcast <- Event{Type: EventTypeConnect, Content: clientID + " connected", Client: clientID}
	}
	mutex.Unlock()

	for {
		var msg Event
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			mutex.Lock()

			if isMonitor {
				delete(monitors, ws)
			} else {
				delete(clients, ws)
				broadcast <- Event{Type: EventTypeDisconnect, Content: clientID + " disconnected", Client: clientID}
			}
			mutex.Unlock()
			break
		}

		if isMonitor && msg.Target != "" {
			// Send message to specific user client
			mutex.Lock()
			for client, id := range clients {
				if id == msg.Target {
					client.WriteJSON(msg)
					break
				}
			}
			mutex.Unlock()
		} else {
			msg.Client = clientID
			broadcast <- msg
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			if msg.Type == EventTypeMessage && clients[client] != msg.Client {
				continue
			}
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		for monitor := range monitors {
			err := monitor.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				monitor.Close()
				delete(monitors, monitor)
			}
		}
		mutex.Unlock()
	}
}
