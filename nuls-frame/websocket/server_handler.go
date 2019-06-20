package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var(
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin:func(r *http.Request) bool{
			return true
		},
	}
)

func NewServerHandler(w http.ResponseWriter , r *http.Request, hub *Hub) {
	var(
		wsConn *websocket.Conn
		err error
		conn *Connection
	)
	// 完成ws协议的握手操作
	// Upgrade:websocket
	if wsConn , err = upgrader.Upgrade(w,r,nil); err != nil {
		return
	}

	conn = &Connection {
		connType:	connTypeIn,
		hub:		hub,
		wsConnect:	wsConn,
		closeChan: 	make(chan byte, 1),
		send:		make(chan []byte, 256),
	}
	conn.hub.register <- conn

	go conn.readPump()
	go conn.writePump()
}