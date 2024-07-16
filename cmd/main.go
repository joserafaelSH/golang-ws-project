package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

//Create a server that can handle ws connections

type ServerConnectionsHandler struct {
	connections      map[*websocket.Conn]bool
	connectionsCount uint8
}

func (s *ServerConnectionsHandler) status() {
	fmt.Printf("Server has %d connections\n", s.connectionsCount)
}

func (s *ServerConnectionsHandler) broadCast(message string) {

	/*
	   for conn, active := range s.connections {

	       if !active {
	           continue
	       }

	       conn.Write([]byte(message))
	   }
	*/

	byteMessage := []byte(message)

	for conn := range s.connections {

		go func(conn *websocket.Conn) {
			if _, err := conn.Write(byteMessage); err != nil {
				fmt.Println("Write error: ", err)
			}
		}(conn)
	}
}

func (s *ServerConnectionsHandler) HandleWsConnections(conn *websocket.Conn) {
	//_, ok := s.connections[conn]
	//if !ok {
	//	s.connections[conn] = true
	//	s.connectionsCount += 1
	//	s.status()
	//}
	s.connections[conn] = true
	s.connectionsCount += 1
	s.status()

	//watch for messages
	s.readLoop(conn)

}

func (s *ServerConnectionsHandler) readLoop(conn *websocket.Conn) {
	buffer := make([]byte, 1024)
	for {

		length, err := conn.Read(buffer)

		if err != nil {
			//some errors like invalid inputs or EOF
			//if you retur, de ReadLoop will be closed
			continue

		}

		msg := buffer[:length]
		conn.Write([]byte("Thank you for the message"))
		fmt.Printf("Client %s say: %s\n", conn.RemoteAddr(), string(msg))
        s.broadCast(string(msg))
	}
}

func CreateConnectionsHandler() *ServerConnectionsHandler {
	return &ServerConnectionsHandler{make(map[*websocket.Conn]bool), 0}
}

func main() {

	handler := CreateConnectionsHandler()

	fmt.Println("Starting webserver on port: 3000")
	http.Handle("/ws", websocket.Handler(handler.HandleWsConnections))
	http.ListenAndServe(":3000", nil)
}
