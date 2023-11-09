package model

type Model struct {
	MsgType     string `json:"message_type"`
	Description string `json:"description"`
	Data        string `json:"data,omitempty"`
}
