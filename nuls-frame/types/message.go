package types

import "encoding/json"

const (
	MessageTypeRequest = "Request"
	MessageTypeResponse = "Response"
)

type Message struct {
	MessageID 		string				`json:"MessageID"`
	Timestamp		string				`json:"Timestamp"`
	TimeZone		string				`json:"TimeZone"`
	MessageType		string				`json:"MessageType"`
	MessageData		json.RawMessage		`json:"MessageData"`
}