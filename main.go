package main

import (
	"flag"
	"log"
	"net"
	"net/url"

	"github.com/gorilla/websocket"
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

	conn, _, err := d.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting WebSocket:", err)
	}

	hilMock := NewHilMock()
	hilMock.SetBackConn(conn)

	hilMock.startIDLE()
}
