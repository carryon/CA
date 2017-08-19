package tables

import (
	"database/sql"
	"time"
)

type Node struct {
	ID         uint64    `json:"f_id"`
	Version    int       `json:"f_version"`
	ChainID    string    `json:"f_chain_id"`
	NodeID     string    `json:"f_node_id"`
	AgentID    string    `json:"f_agent_id"`
	ConfigFile string    `json:"f_config_file"`
	Config     string    `json:"f_config"`
	Status     string    `json:"f_status"`
	Height     int       `json:"f_height"`
	Addr       string    `json:"f_addr"`
	Created    time.Time `json:"f_created_at"`
	Updated    time.Time `json:"f_updated_at"`
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) Insert(tx *sql.Tx) error {

	return nil
}

func (n *Node) Update(tx *sql.Tx) error {

	return nil
}

func (n *Node) Delete(tx *sql.Tx, condition string) error {

	return nil
}
