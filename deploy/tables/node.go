package tables

import (
	"database/sql"
	"fmt"
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
	//todo
	return nil
}

func (n *Node) UpdateHeight(tx *sql.Tx) error {
	res, err := tx.Exec("update t_node set f_height=?, f_status=?, f_updated_at=? where f_chain_id=? and f_node_id=?", n.Height, n.Status, n.Updated, n.ChainID, n.NodeID)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (n *Node) UpdateAllConfig(tx *sql.Tx) error {
	res, err := tx.Exec("update t_node set f_config=? where f_chain_id=? ", n.Config, n.ChainID)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found ")
	}
	return nil
}

func (n *Node) Delete(tx *sql.Tx, condition string) error {
	//todo
	return nil
}
