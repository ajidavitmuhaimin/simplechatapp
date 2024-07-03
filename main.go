package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

type Message struct {
    Username string `json:"username"`
    Message  string `json:"message"`
}

func main() {
    http.HandleFunc("/", handleHome)
    http.HandleFunc("/ws", handleConnections)
    go handleMessages()

    fmt.Println("Server started on :1395")
    err := http.ListenAndServe(":1395", nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer ws.Close()

    clients[ws] = true

    for {
        var msg Message
        err := ws.ReadJSON(&msg)
        if err != nil {
            fmt.Println("Error reading JSON:", err)
            delete(clients, ws)
            break
        }
        broadcast <- msg
    }
}

func handleMessages() {
    for {
        msg := <-broadcast
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                fmt.Println("Error writing JSON:", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}
