package main

type Notification struct {
	MsgType     string `json:"message_type" bson:"message_type"`
	Description string `json:"description" bson:"description"`
	Data        string `json:"data,omitempty" bson:"data,omitempty"`
}
