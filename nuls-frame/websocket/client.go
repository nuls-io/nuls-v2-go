package websocket

import (
	"flag"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func NewClient(addr, uri string) (*Connection, error) {
	var addrStr = flag.String("addr", addr, "http service address")

	u := url.URL{Scheme: "ws", Host: *addrStr, Path: uri}
	log.Printf("connecting to %s", u.String())

	wsConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	conn := &Connection{
		connType:	connTypeOut,
		hub:		GetHub(),
		wsConnect:	wsConn,
		closeChan: 	make(chan byte, 1),
		send:		make(chan []byte, 256),
	}

	conn.hub.register <- conn

	go conn.readPump()
	go conn.writePump()

	return conn, nil
}
