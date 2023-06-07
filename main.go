package main

import (
	"flag"
	"net"
	"net/url"

	"github.com/gorilla/websocket"
	trace "github.com/rs/zerolog/log"
)

var traceLevel = flag.String("trace", "info", "set the trace level (\"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\")")
var traceFile = flag.String("log", "trace.json", "set the trace log file")

func main() {
	traceFile := initTrace(*traceLevel, *traceFile)
	defer traceFile.Close()

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8010", Path: "/backend"}

	d := websocket.Dialer{
		NetDialContext: (&net.Dialer{LocalAddr: &net.TCPAddr{IP: net.ParseIP("127.0.0.2")}}).DialContext,
	}

	for {
		conn, _, err := d.Dial(u.String(), nil)
		if err != nil {
			trace.Error().Err(err).Msg("Error connecting WebSocket")
			return
		}

		hilMock := NewHilMock()
		hilMock.SetBackConn(conn)

		hilMock.startIDLE()
	}
}
