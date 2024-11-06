package models

type ChatMessage struct {
	User string `json:"user"`
	Msg  string `json:"msg"`
}
