package tables

import (
	"database/sql"
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
		// if err := rows.Scan(
		// 	&node.ID,
		// 	&node.Version,
		// 	&node.ChainID,
		// 	&node.NodeID,
		// 	&node.ConfigFile,
		// 	&node.Config,
		// 	&node.Status,
		// 	&node.Height,
		// 	&node.Addr,
		// 	&node.Created,
		// 	&node.Updated); err != nil {
		// 	return nil, err
		// }
		if err := rows.Scan(&node); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func QueryAllAgent(db *sql.DB) ([]*Agent, error) {
	sql := "select *from t_agent"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*Agent
	for rows.Next() {
		agent := NewAgent()
		// if err := rows.Scan(
		// 	&agent.ID,
		// 	&agent.AgentID,
		// 	&agent.Addr,
		// 	&agent.Created,
		// ); err != nil {
		// 	return nil, err
		// }
		if err := rows.Scan(&agent); err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}

	return agents, nil
}
