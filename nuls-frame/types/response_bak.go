package types
//
//import "time"
//
//type Response map[string]interface{}
//
//func Convert(resMap *map[string]interface{}) *Response{
//	var response Response = *resMap
//	return &response
//}
//
//func NewResponse(requestID string, responseStatus string, responseComment string, responseMaxSize string, responseData *map[string]interface{}) *Response {
//	response := &Response{
//		"RequestID":					requestID,
//		"ResponseProcessingTime":		time.Now().String(),
//		"ResponseStatus":				responseStatus,
//		"ResponseComment": 				responseComment,
//		"ResponseMaxSize":				responseMaxSize,
//		"ResponseData":					responseData,
//	}
//	return response
//}
//
//func NewCmdNotFindResponse(messageId string, cmdName string) *Response {
//	response := &Response{
//		"RequestID":					messageId,
//		"ResponseProcessingTime":		time.Now().String(),
//		"ResponseStatus":				"1",
//		"ResponseComment": 				"",
//		"ResponseMaxSize":				"",
//	}
//	return response
//}
//
//func (response *Response) GetRequestID() string {
//	return (*response)["RequestID"].(string)
//}
//
//func (response *Response) GetResponseProcessingTime() string {
//	return (*response)["ResponseProcessingTime"].(string)
//}
//
//func (response *Response) GetResponseStatus() string {
//	return (*response)["ResponseStatus"].(string)
//}
//
//func (response *Response) GetResponseComment() string {
//	return (*response)["ResponseComment"].(string)
//}
//
//func (response *Response) GetResponseMaxSize() string {
//	return (*response)["ResponseMaxSize"].(string)
//}
//
//func (response *Response) GetResponseData() *map[string]interface{} {
//	return (*response)["ResponseData"].(*map[string]interface{})
//}
//
//
//func (response *Response) SetRequestID(requestId string) {
//	(*response)["RequestID"] = requestId
//}