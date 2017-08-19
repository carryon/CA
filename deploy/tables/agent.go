package tables

import (
	"database/sql"
	"time"
)

type Agent struct {
	ID      uint64    `json:"f_id"`
	AgentID string    `json:"f_agent_id"`
	Addr    string    `jso:"f_addr"`
	Created time.Time `json:"f_create_at"`
}

func NewAgent() *Agent {
	return &Agent{}
}

func (a *Agent) Insert(tx *sql.Tx) error {

	return nil
}

func (a *Agent) Update(tx *sql.Tx) error {

	return nil
}

func (a *Agent) Delete(tx *sql.Tx, condition string) error {

	return nil
}
