package model

// Model представляет структуру данных для модели уведомления.
type Model struct {
	MsgType     string `json:"message_type" bson:"message_type"`
	Description string `json:"description" bson:"description"`
	Date        string `json:"date" bson:"date"`
}
