package server

import (
	"github.com/bocheninc/CA/deploy/components/utils"
)

const (
	ChainChangeCfgMsg = iota + 106
	ChainNodeStatusMsg
)

type ChangeCfg struct {
	ChainID string
	Message *Message
}

// Message represents the message transfer in msg-net
type Message struct {
	Cmd     uint8
	Payload []byte
}

// Serialize message to bytes
func (m *Message) Serialize() []byte {
	return utils.Serialize(*m)
}

// Deserialize bytes to message
func (m *Message) Deserialize(data []byte) {
	utils.Deserialize(data, m)
}
