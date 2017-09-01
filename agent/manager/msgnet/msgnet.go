package msgnet

import (
	"sync"

	"github.com/bocheninc/CA/agent/types"
)

type MsgnetInfo struct {
	sync.RWMutex
	MsgnetID string
	Config   *types.MsgnetConfig
	Status   bool
	Version  string
}

func NewMsgnetInfo(ID, version string, config *types.MsgnetConfig) *MsgnetInfo {
	return &MsgnetInfo{MsgnetID: ID,
		Version: version,
		Config:  config,
		Status:  false}
}

func (m *MsgnetInfo) CheckVersion(version string) bool {
	if m.Version == version {
		return true
	}
	return false
}

func (m *MsgnetInfo) Start() error {
	if !m.Status {

		m.Status = true
	}

	return nil
}

func (m *MsgnetInfo) Stop() error {
	return nil
}
