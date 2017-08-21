package manager

import (
	"sync"

	"github.com/bocheninc/CA/Agent-Server/manager/msgnet"
	"github.com/bocheninc/CA/Agent-Server/manager/node"
)

type Manager struct {
	sync.RWMutex
	NodeList   map[string]*node.NodeInfo
	MsgnetList map[string]*msgnet.MsgnetInfo
}

func NewManager() *Manager {
	return &Manager{}
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

func (m *Manager) UpdateNodes(nodes []*node.NodeInfo) error {
	m.Lock()
	defer m.Unlock()

	//add and update version
	for _, v := range nodes {
		if node, ok := m.NodeList[v.NodeID]; ok {
			if !node.CheckVersion(node.Version) {
				m.NodeList[v.NodeID] = v
			}
		} else {
			m.NodeList[v.NodeID] = v
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
