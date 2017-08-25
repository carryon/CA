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
		tmpNodeList[v.AgentID] = append(tmpNodeList[v.AgentID], v)
		if _, ok := l.NodeList[v.AgentID]; ok {
			for _, node := range l.NodeList[v.AgentID] {
				if v.NodeID == node.NodeID {
					// node exist but node.config change
					if v.Config != node.Config {
						//make msg

						msg := new(Message)
						msg.Cmd = ChainChangeCfgMsg
						msg.Payload = []byte(v.Config)
						l.msgChan <- &ChangeCfg{ChainID: v.ChainID, Message: msg}

						//update other node config
						tx, _ := l.Db.Begin()
						log.Debug("begin...", v.NodeID, v.Config, node.NodeID, node.Config)
						if err := v.UpdateAllConfig(tx); err != nil {
							log.Errorf(" chainID: %s, nodeID: %s update config %s, err: %v", node.ChainID, node.NodeID, node.Config, err)
							tx.Rollback()
						}
						tx.Commit()
					}

				}

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

func (l *List) nodeIsExist(node *tables.Node, array []*tables.Node) bool {
	for _, v := range array {
		if node.NodeID == v.NodeID {
			return true
		}
	}
	return false
}
