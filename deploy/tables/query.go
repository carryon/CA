package tables

import (
	"database/sql"
	"fmt"
	"strconv"
)

func QueryAllNode(db *sql.DB) ([]*Node, error) {

	sql := "select * from t_node"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []*Node
	for rows.Next() {
		node := NewNode()
		if err := rows.Scan(
			&node.ID,
			&node.Version,
			&node.ChainID,
			&node.NodeID,
			&node.AgentID,
			&node.ConfigFile,
			&node.Config,
			&node.Status,
			&node.Height,
			&node.Addr,
			&node.Created,
			&node.Updated); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func QueryAllAgent(db *sql.DB) ([]*Agent, error) {
	sql := "select * from t_agent"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*Agent
	for rows.Next() {
		agent := NewAgent()
		if err := rows.Scan(
			&agent.ID,
			&agent.AgentID,
			&agent.Addr,
			&agent.Created,
		); err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}

	return agents, nil
}

func QueryAllTps(db *sql.DB) (int, error) {

	type tps struct {
		ChainID string `json:"f_chain_id"`
		Tps     string `json:"f_status"`
	}

	sql := "SELECT f_chain_id, f_status FROM t_node WHERE f_updated_at >= NOW()-INTERVAL 60 SECOND"
	rows, err := db.Query(sql)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var allTps int
	for rows.Next() {
		tps := new(tps)
		if err := rows.Scan(
			&tps.ChainID,
			&tps.Tps,
		); err != nil {
			return 0, err
		}
		t, err := strconv.Atoi(tps.Tps)
		if err != nil {
			return 0, err
		}
		allTps += t
	}

	return allTps / 4, nil
}

func QueryCertByID(db *sql.DB, chainID, nodeID string) ([]*Cert, error) {
	sql := fmt.Sprintf("select f_id, f_chain_id, f_node_id, f_crt, f_publicKey, f_root_crt, f_created_at from t_cert where f_chain_id='%s' and f_node_id='%s'", chainID, nodeID)
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var certs []*Cert
	for rows.Next() {
		cert := new(Cert)
		if err := rows.Scan(
			&cert.ID,
			&cert.ChainID,
			&cert.NodeID,
			&cert.Crt,
			&cert.PublicKey,
			&cert.RootCrt,
			&cert.Created,
		); err != nil {
			return nil, err
		}

		certs = append(certs, cert)
	}
	return certs, nil
}
