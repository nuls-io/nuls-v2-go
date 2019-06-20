package websocket

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/json-iterator/go"

	"github.com/nuls-io/go-nuls/nuls-frame/types"
	"github.com/nuls-io/go-nuls/nuls-frame/util"
)
var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 20971520

	connTypeIn = 0

	connTypeOut = 1
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
	autoInc = util.New(1, 1)

	//RequestMaps = make(map[int]chan *Message)
	RequestMaps = sync.Map{}
)

type Connection struct{
	connType	int
	hub			*Hub
	wsConnect 	*websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
	closeChan 	chan byte
	mutex 		sync.Mutex  // 对closeChan关闭上锁
	isClosed 	bool  // 防止closeChan被关闭多次
}

func (conn *Connection) SendRequestAndGetResponse(request *types.Request) (*types.Response, error) {

	t := time.Now()
	id, waitChan, error := conn.SendRequest(request)

	if error != nil {
		return nil ,error
	}
	util.T5 += time.Now().Sub(t)

	now := time.Now()
	//wait response
	res := <- waitChan
	util.T1 += time.Now().Sub(now)

	RequestMaps.Delete(id)

	return res, nil
}

func (conn *Connection) SendRequest(request *types.Request) (string, chan *types.Response, error) {

	id := strconv.Itoa(autoInc.Id())

	messageData, _ := json.Marshal(request)
	message := &types.Message{
		MessageType:		types.MessageTypeRequest,
		MessageID:			id,
		Timestamp:			time.Now().String(),
		MessageData:		messageData,
	}
	bytes, error := json.Marshal(message)
	if error != nil {
		return "", nil, error
	}

	waitChan := make(chan *types.Response)
	RequestMaps.Store(id ,waitChan)

	conn.send <- bytes
	return id, waitChan, nil
}

func (conn *Connection) SendResponse(response *types.Response) error {
	id := strconv.Itoa(autoInc.Id())
	messageData, _ := json.Marshal(response)
	message := &types.Message{
		MessageType:		types.MessageTypeResponse,
		MessageID:			id,
		Timestamp:			time.Now().String(),
		MessageData:		messageData,
	}

	bytes, error := json.Marshal(message)
	if error != nil {
		return error
	}
	conn.send <- bytes
	return nil
}

func (conn *Connection) Close() {
	// 线程安全，可多次调用
	conn.wsConnect.Close()
	// 利用标记，让closeChan只关闭一次
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Connection) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.wsConnect.Close()
	}()
	c.wsConnect.SetReadLimit(maxMessageSize)
	c.wsConnect.SetReadDeadline(time.Now().Add(pongWait))
	c.wsConnect.SetPongHandler(func(string) error { c.wsConnect.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.wsConnect.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		//msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		now := time.Now()
		message := &types.Message{}
		error := json.Unmarshal(msg, message)
		if error != nil {
			log.Println("error data format : ", string(msg))
			continue
		}
		util.T4 += time.Now().Sub(now)
		c.hub.handler.OnMessage(message, c)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.wsConnect.Close()
	}()
	for {
		select {
		case message, ok := <- c.send:
			c.wsConnect.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.wsConnect.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.wsConnect.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//	w.Write(newline)
			//	w.Write(<-c.send)
			//}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.wsConnect.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.wsConnect.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}