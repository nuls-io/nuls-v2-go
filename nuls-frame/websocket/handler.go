package websocket

import (
	"log"

	"github.com/nuls-io/go-nuls/nuls-frame/types"
)

type Handler interface {
	OnMessage(msg *types.Message, conn *Connection) error
}

type DefaultHandler struct {
}

func (handler *DefaultHandler) OnMessage(msg *types.Message, conn *Connection) error {
	if conn.connType == connTypeIn {
		log.Println("Server receive message : %V", msg)

		replay := make(map[string]interface{})
		replay["msg"] = "this is Server replay of message [" + msg.MessageID + "]"
		response := types.NewResponse(msg.MessageID, "1", "test", "0",	&replay)

		conn.SendResponse(response)
	} else if conn.connType == connTypeOut {
		log.Println("Client receive message : %V", msg)
	}
	return nil
}