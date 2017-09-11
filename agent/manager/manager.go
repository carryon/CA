package manager

import (
	"sync"

	"github.com/bocheninc/CA/agent/config"
	"github.com/bocheninc/CA/agent/log"
	"github.com/bocheninc/CA/agent/manager/msgnet"
	"github.com/bocheninc/CA/agent/manager/node"
	"github.com/bocheninc/CA/agent/request"
)

type Manager struct {
	sync.RWMutex
	req        *request.Request
	NodeList   map[string]*node.NodeInfo
	MsgnetList map[string]*msgnet.MsgnetInfo
}

func NewManager() *Manager {
	return &Manager{
		req:        request.NewRequest(),
		NodeList:   make(map[string]*node.NodeInfo),
		MsgnetList: make(map[string]*msgnet.MsgnetInfo),
	}
}

func (m *Manager) StartNodes() error {
	m.Lock()
	defer m.Unlock()
	for _, node := range m.NodeList {
		if err := node.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) StartMsgnets() error {
	m.Lock()
	defer m.Unlock()
	for _, msgnet := range m.MsgnetList {
		if err := msgnet.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) UpdateNodes() error {
	m.Lock()
	defer m.Unlock()

	nodes, err := m.req.GetLcndConfig(config.Cfg.ID)
	if err != nil {
		log.Error("Get nodes config err: ", err)
		return err
	}

	//add and update version
	for _, v := range nodes {
		if node, ok := m.NodeList[v.NodeID]; ok {
			if !node.CheckVersion(node.Version) {
				m.NodeList[v.NodeID] = v
				cert, err := m.req.GetCrt(v.Config.Blockchain.ChainID, v.Config.Blockchain.NodeID)
				if err != nil {
					return err
				}
				m.NodeList[v.NodeID].Cert = cert
			}
		} else {
			m.NodeList[v.NodeID] = v
			cert, err := m.req.GetCrt(v.Config.Blockchain.ChainID, v.Config.Blockchain.NodeID)
			if err != nil {
				return err
			}
			m.NodeList[v.NodeID].Cert = cert
		}
	}

	//delete
	for id, node := range m.NodeList {
		var flag = false
		for _, v := range nodes {
			if v.NodeID == id {
				flag = true
			}
		}

		if !flag {
			if err := node.Stop(); err != nil {
				return err
			}
			delete(m.NodeList, id)
		}
	}
	return nil
}

func (m *Manager) UpdateMsgnets(msgnets []*msgnet.MsgnetInfo) error {
	m.Lock()
	defer m.Unlock()

	//add and update version
	for _, v := range msgnets {
		if msgnet, ok := m.MsgnetList[v.MsgnetID]; ok {
			if !msgnet.CheckVersion(msgnet.Version) {
				m.MsgnetList[v.MsgnetID] = v
			}
		} else {
			m.MsgnetList[v.MsgnetID] = v
		}
	}

	//delete
	for id, msgnet := range m.MsgnetList {
		var flag = false
		for _, v := range msgnets {
			if v.MsgnetID == id {
				flag = true
			}
		}

		if !flag {
			if err := msgnet.Stop(); err != nil {
				return err
			}
			delete(m.MsgnetList, id)
		}
	}
	return nil
}
