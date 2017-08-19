package server

import (
	"database/sql"
	"sync"

	"github.com/bocheninc/CA/deploy/components/log"

	"github.com/bocheninc/CA/deploy/tables"
)

type List struct {
	db *sql.DB
	sync.Mutex
	NodeList  map[string][]*tables.Node
	AgentList map[string]*tables.Agent
}

func NewList(db *sql.DB) *List {
	return &List{
		NodeList:  make(map[string][]*tables.Node),
		AgentList: make(map[string]*tables.Agent),
		db:        db}
}

func (l *List) UpdateNodeList() {
	l.Lock()
	defer l.Unlock()
	tmpNodeList := make(map[string][]*tables.Node)

	agents, err := tables.QueryAllAgent(l.db)
	if err != nil {
		log.Error("query all agent error: ", err)
	}

	for _, v := range agents {
		l.AgentList[v.AgentID] = v
	}

	nodes, err := tables.QueryAllNode(l.db)
	if err != nil {
		log.Error("query all node error: ", err)
	}

	for _, v := range nodes {
		tmpNodeList[v.AgentID] = append(tmpNodeList[v.AgentID], v)
		if _, ok := l.NodeList[v.AgentID]; ok {
			if l.nodeIsExist(v, l.NodeList[v.AgentID]) {
				//Todo node exist but node.config change
			}
		}
	}

	l.NodeList = tmpNodeList
}

func (l *List) nodeIsExist(node *tables.Node, array []*tables.Node) bool {
	for _, v := range array {
		if node.NodeID == v.NodeID {
			return true
		}
	}
	return false
}
