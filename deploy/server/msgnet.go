package server

import (
	msgnet "github.com/bocheninc/msg-net/peer"
)

// MsgHandler handles the message of the msg-net
type MsgHandler func(src string, dst string, payload, sig []byte) error

type MsgNet struct {
	peer *msgnet.Peer
}

// NewMsgnet start client msg-net service and returns a msg-net peer
func NewMsgnet(id string, routeAddress []string, fn MsgHandler, logOutPath string) *MsgNet {
	// msg-net services
	msgnet.SetLogOut(logOutPath)
	msgnet.NewPeer(id, routeAddress, fn)
	return &MsgNet{peer: msgnet.NewPeer(id, routeAddress, fn)}

}

func (m *MsgNet) Start() {
	m.peer.Start()
}
