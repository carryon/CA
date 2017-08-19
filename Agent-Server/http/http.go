package http

import "encoding/json"

type Response struct {
	Result []string    `json:"result"`
	Error  interface{} `json:"error"`
	ID     uint        `json:"id"`
}

type Request struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	ID     uint             `json:"id"`
}
