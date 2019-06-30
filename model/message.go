package model

// Message message style
type Message struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"message"`
}
