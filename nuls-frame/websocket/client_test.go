package websocket

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/nuls-io/go-nuls/nuls-frame/types"
)

func TestOnMessage(t *testing.T) {

	startupServer()

	addr := "127.0.0.1:7777"
	uri := "/ws"

	conn, error := NewClient(addr, uri)

	if error != nil {
		log.Fatal(error)
	}

	defer conn.Close()

	for i := 0 ; i < 10 ; i ++ {
		params := map[string]string {
			"params" : "this is a test of websocket , the index is " + strconv.Itoa(i),
			"version" : "1.0",
			"count" : "" + strconv.Itoa(i),
		}
		methods := map[string]map[string]string {
			"sum" : params,
		}

		request := &types.Request{
			RequestAck : "0",
			SubscriptionEventCounter : "0",
			SubscriptionPeriod : "0",
			SubscriptionRange : "0",
			ResponseMaxSize : "1000",
			RequestMethods : methods,
		}

		res, _, error := conn.SendRequest(request)
		if error != nil {
			log.Println(" error : ", error)
		}

		if error != nil {

			log.Println("receive response : " , res)
		}
	}

	time.Sleep(time.Second)
}

func startupServer() {
	addr := "127.0.0.1:7777"
	uri := "/ws"
	handler := &DefaultHandler{}

	server := &Server{
		Addr:		addr,
		Uri:		uri,
		Handler:	handler,
	}

	server.Startup()
}
