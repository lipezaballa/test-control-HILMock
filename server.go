package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	trace "github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	handleConn func(conn *websocket.Conn)
}

func NewServer() Server {
	return Server{
		handleConn: func(conn *websocket.Conn) {},
	}

}

func (server *Server) SetConnHandler(handler func(conn *websocket.Conn)) { //TODO: add as prop handler
	server.handleConn = handler
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		trace.Error().Err(err).Msg("Error upgrading connection")
		return
	}
	server.handleConn(conn)
}
