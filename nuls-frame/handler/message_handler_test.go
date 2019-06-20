package handler

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/nuls-io/go-nuls/nuls-frame/goroutine"
	"github.com/nuls-io/go-nuls/nuls-frame/types"
	"github.com/nuls-io/go-nuls/nuls-frame/util"
	"github.com/nuls-io/go-nuls/nuls-frame/websocket"
)

func TestOnMessage(t *testing.T) {
	addr := "127.0.0.1:7777"
	//remoteAddr := "192.168.1.202:8887"
	uri := "/ws"

	server := &websocket.Server{
		Addr:		addr,
		Uri:		uri,
		Handler:	&TextMessageHandler{},
	}

	server.Startup()

	conn, error := websocket.NewClient(addr, uri)

	if error != nil {
		log.Fatal(error)
	}

	defer conn.Close()

	now := time.Now()

	for i := 0 ; i < 100000 ; i ++ {

		params := map[string]interface{} {
			"params" : "this a test of websocket , the index is " + strconv.Itoa(i),
			"version" : "1.0",
			"count" : strconv.Itoa(i),
		}
		methods := &map[string]map[string]interface{} {
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


		res, error := conn.SendRequestAndGetResponse(request)
		if error != nil {
			log.Println(" error : ", error)
		}

		if error != nil {
			log.Println("receive response : " , res)
		}
	}

	diff := time.Now().Sub(now)

	log.Println("use time : ", diff)

	util.PrintTime()
}


func TestOnMessage_Multi(t *testing.T) {
	addr := "127.0.0.1:7777"
	uri := "/ws"

	server := &websocket.Server{
		Addr:		addr,
		Uri:		uri,
		Handler:	&TextMessageHandler{},
	}

	server.Startup()

	conn, error := websocket.NewClient(addr, uri)

	if error != nil {
		log.Fatal(error)
	}

	defer conn.Close()

	now := time.Now()

	pool := goroutine.NewPool(8)
	pool.Run()

	total := 100000
	wait := make(chan *types.Response, total)

	for i := 0 ; i < total ; i ++ {
		pool.Execute(goroutine.NewTask(func() {
			params := map[string]interface{} {
				"params" : "this a test of websocket , the index is " + strconv.Itoa(i),
				"version" : "1.0",
				"count" : "" + strconv.Itoa(i),
			}
			methods := &map[string]map[string]interface{} {
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
			res, error := conn.SendRequestAndGetResponse(request)
			if error != nil {
				log.Println(" error : ", error)
			}
			if error != nil {
				log.Println("receive response : " , res)
			}
			wait <- res
		}))
	}

	for i := 0 ; i < total ; i ++ {
		<- wait
	}

	diff := time.Now().Sub(now)

	log.Println("use time : ", diff)

	util.PrintTime()
}
