package types

import "time"

type Response struct {
	RequestID				string						`json:"RequestID"`
	ResponseProcessingTime	string						`json:"ResponseProcessingTime"`
	ResponseStatus			string						`json:"ResponseStatus"`
	ResponseComment			string						`json:"ResponseComment"`
	ResponseMaxSize			string						`json:"ResponseMaxSize"`
	ResponseData			*map[string]interface{}		`json:"ResponseData"`
}

func NewResponse(requestID string, responseStatus string, responseComment string, responseMaxSize string, responseData *map[string]interface{}) *Response {
	response := &Response{
		RequestID:						requestID,
		ResponseProcessingTime:			time.Now().String(),
		ResponseStatus:					responseStatus,
		ResponseComment: 				responseComment,
		ResponseMaxSize:				responseMaxSize,
		ResponseData:					responseData,
	}
	return response
}

func NewCmdNotFindResponse(messageId string, cmdName string) *Response {
	response := &Response{
		RequestID:						messageId,
		ResponseProcessingTime:			time.Now().Format("2006-01-02 15:04:05"),
		ResponseStatus:					"1",
		ResponseComment: 				"",
		ResponseMaxSize:				"",
	}
	return response
}

func NewErrorResponse(messageId string, cmdName string, error error) *Response {
	response := &Response{
		RequestID:					messageId,
		ResponseProcessingTime:		time.Now().String(),
		ResponseStatus:				"1",
		ResponseComment: 				"",
		ResponseMaxSize:				"",
	}
	return response
}
