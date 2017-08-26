package server

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/bocheninc/CA/deploy/components/log"
	"github.com/bocheninc/CA/deploy/tables"
)

type List struct {
	Db *sql.DB
	sync.Mutex
	NodeList  map[string][]*tables.Node
	AgentList map[string]*tables.Agent
	msgChan   chan *ChangeCfg
}

func NewList(db *sql.DB) *List {
	return &List{
		NodeList:  make(map[string][]*tables.Node),
		AgentList: make(map[string]*tables.Agent),
		msgChan:   make(chan *ChangeCfg, 100),
		Db:        db}
}

func (l *List) UpdateNodeList() {
	l.Lock()
	defer l.Unlock()

	type nodeConfig struct {
		update bool
		config string
	}

	updateConfig := make(map[string]*nodeConfig)
	tmpNodeList := make(map[string][]*tables.Node)

	agents, err := tables.QueryAllAgent(l.Db)
	if err != nil {
		log.Error("query all agent error: ", err)
	}

	for _, v := range agents {
		l.AgentList[v.AgentID] = v
	}

	nodes, err := tables.QueryAllNode(l.Db)
	if err != nil {
		log.Error("query all node error: ", err)
	}
	for _, v := range nodes {
		log.Debug("query node info", v.NodeID, " config: ", v.Config)
	}

	for _, v := range nodes {
		tmpNodeList[v.AgentID] = append(tmpNodeList[v.AgentID], v)
		//initialization
		if _, ok := updateConfig[v.AgentID]; !ok {
			updateConfig[v.AgentID] = &nodeConfig{update: false}
		}

		if _, ok := l.NodeList[v.AgentID]; ok {
			if !updateConfig[v.AgentID].update {
				if l.updateAllConfig(v, l.NodeList[v.AgentID]) {
					updateConfig[v.AgentID] = &nodeConfig{update: true, config: v.Config}
				}
			}
		}
	}

	for k, v := range updateConfig {
		if v.update {
			for _, node := range tmpNodeList[k] {
				node.Config = v.config
			}
		}
	}

	l.NodeList = tmpNodeList
}

func (l *List) GetNodeByNodeID(chainID, nodeID string) (*tables.Node, error) {
	l.Lock()
	defer l.Unlock()

	for _, v := range l.NodeList {
		for _, node := range v {
			if chainID == node.ChainID && nodeID == node.NodeID {
				return node, nil
			}
		}
	}

	return nil, fmt.Errorf("not found node by chainID: %s and nodeID %s .", chainID, nodeID)
}

func (l *List) MessageChan() <-chan *ChangeCfg {
	return l.msgChan
}

func (l *List) updateAllConfig(node *tables.Node, array []*tables.Node) bool {
	for _, v := range array {
		if node.NodeID == v.NodeID {
			if node.Config != v.Config {
				log.Debug("begin...", v.NodeID, " config: ", v.Config, " nodeID: ", node.NodeID, " node.config: ", node.Config)
				//make msg
				msg := new(Message)
				msg.Cmd = ChainChangeCfgMsg
				msg.Payload = []byte(node.Config)
				l.msgChan <- &ChangeCfg{ChainID: v.ChainID, Message: msg}

				//update other node config
				tx, _ := l.Db.Begin()
				if err := node.UpdateAllConfig(tx); err != nil {
					log.Errorf(" chainID: %s, nodeID: %s update config %s, err: %v", node.ChainID, node.NodeID, node.Config, err)
					tx.Rollback()
				}
				tx.Commit()
				return true
			}
		}
	}
	return false
}
