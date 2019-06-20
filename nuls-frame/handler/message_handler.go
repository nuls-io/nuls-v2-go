package handler

import (
	"errors"
	"log"
	"time"

	"github.com/json-iterator/go"

	"github.com/nuls-io/go-nuls/nuls-frame/cmd"
	"github.com/nuls-io/go-nuls/nuls-frame/goroutine"
	"github.com/nuls-io/go-nuls/nuls-frame/types"
	"github.com/nuls-io/go-nuls/nuls-frame/websocket"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type TextMessageHandler struct {
}

var(
	requestPool *goroutine.Pool
	responsePool *goroutine.Pool
)

func init() {
	requestPool = goroutine.NewPool(4)
	go requestPool.Run()

	responsePool = goroutine.NewPool(4)
	go responsePool.Run()
}

func (handler *TextMessageHandler) OnMessage(msg *types.Message, conn *websocket.Connection) error {

	if msg.MessageType == types.MessageTypeRequest {
		f := func() {
			error := handlerRequest(msg, conn)
			if error != nil {
				log.Println(error)
			}
		}
		requestPool.Execute(goroutine.NewTask(f))
	} else if msg.MessageType == types.MessageTypeResponse {
		f := func() {
			error := handlerResponse(msg, conn)
			if error != nil {
				log.Println(error)
			}
		}
		responsePool.Execute(goroutine.NewTask(f))
	} else {
		log.Println("unknown message")
	}
	return nil
}

func handlerResponse(msg *types.Message, conn *websocket.Connection) error {
	response := &types.Response{}
		json.Unmarshal(msg.MessageData, response)
		//if !ok {
		//	errMsg := "unrecognizable data format"
		//	return errors.New(errMsg)
		//}
		//response := types.Convert(resMap)
		requestId := response.RequestID
		if requestId == "" {
			errMsg := "error! receive a response message, but not found request id"
			return errors.New(errMsg)
		}
		for {
			retryCount := 0
			waitChan, ok := websocket.RequestMaps.Load(requestId)
			if ok && waitChan != nil {
				chanValue,success:= waitChan.(chan *types.Response)
				if success {
					chanValue <- response
				} else {
					errMsg := "error , the value is not chan message"
					return errors.New(errMsg)
				}
				break
			} else {
				if retryCount > 1000 {
					break
				}
				time.Sleep(time.Millisecond)
				retryCount++
			}
		}
	return nil
}

func handlerRequest(msg *types.Message, conn *websocket.Connection) error {

	// write log ...

	//log.Println(" msg : ", msg)
	request := &types.Request{}
	json.Unmarshal(msg.MessageData, request)

	// get process function
	//messageData, ok := msg.MessageData.(map[string]interface{})
	//if !ok {
	//	errMsg := "unknown message format"
	//	return errors.New(errMsg)
	//}
	//log.Println("request: ", request)

	for k, v := range *request.RequestMethods {
		cmd := getCmdFunction(k)
		//log.Println("cmd :", cmd)

		if cmd == nil {
			// 提示，不存在的cmd
			// waring , Non-existent CMD
			conn.SendResponse(types.NewCmdNotFindResponse(msg.MessageID, k))
			continue
		}
		res, err := (*cmd).Do(&v)
		if err != nil {
			// 日志记录 ?
			conn.SendResponse(types.NewErrorResponse(msg.MessageID, k, err))
			continue
		}
		res.RequestID = msg.MessageID
		//res.ResponseStatus = "0"
		//res.ResponseMaxSize = "0"
		//res.ResponseProcessingTime = time.Now().Format("2006-01-02 15:04:05")

		conn.SendResponse(res)
	}
	
	return nil
}

func getCmdFunction(name string) *cmd.Cmd {
	return cmd.GetCmd(name)
}