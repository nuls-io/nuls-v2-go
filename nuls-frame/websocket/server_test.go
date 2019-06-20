package websocket

import (
	"github.com/nuls-io/go-nuls/nuls-frame/types"
	"log"
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	addr := "127.0.0.1:7777"
	uri := "/ws"
	handler := &DefaultHandler{}

	server := &Server{
		Addr:		addr,
		Uri:		uri,
		Handler:	handler,
	}

	server.Startup()

	conn, error := NewClient(addr, uri)

	if error != nil {
		log.Fatal(error)
	}

	defer conn.Close()
	for i := 0 ; i < 10 ; i ++ {
		methods := map[string]map[string]string {}
		request := &types.Request{
			RequestAck : "0",
			SubscriptionEventCounter : "0",
			SubscriptionPeriod : "0",
			SubscriptionRange : "0",
			ResponseMaxSize : "1000",
			RequestMethods : methods,
		}
		conn.SendRequest(request)
	}

	time.Sleep(1 * time.Second)
}
