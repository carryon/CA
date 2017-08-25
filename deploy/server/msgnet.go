package server

import (
	"github.com/bocheninc/CA/deploy/components/log"
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
	if len(routeAddress) < 0 {
		return nil
	}

	msgnet.SetLogOut(logOutPath)
	msgnet.NewPeer(id, routeAddress, fn)
	return &MsgNet{peer: msgnet.NewPeer(id, routeAddress, fn)}

}

func (m *MsgNet) Start() {
	m.peer.Start()
	log.Debug("Msg-net Service Start ...")
}

func (m *MsgNet) SendMsgnetMessage(src, dst string, msg *Message) bool {
	//todo signature
	log.Debugf("==============send data=========== cmd : %v, payload: %v ", msg.Cmd, string(msg.Payload))
	return m.peer.Send(dst, msg.Serialize(), []byte("deploy"))

}
