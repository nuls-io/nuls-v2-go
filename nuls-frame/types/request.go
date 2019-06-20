package types

type Request struct {
	RequestAck						string						`json:"RequestAck"`
	SubscriptionEventCounter		string						`json:"SubscriptionEventCounter"`
	SubscriptionPeriod				string						`json:"SubscriptionPeriod"`
	SubscriptionRange				string						`json:"SubscriptionRange"`
	ResponseMaxSize					string						`json:"ResponseMaxSize"`
	RequestMethods					*map[string]map[string]interface{}		`json:"RequestMethods"`
}
