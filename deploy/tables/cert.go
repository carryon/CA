package tables

import (
	"database/sql"
	"fmt"
	"time"
)

type Cert struct {
	ID        uint64    `json:"f_id"`
	ChainID   string    `json:"f_chain_id"`
	NodeID    string    `json:"f_node_id"`
	Crt       string    `json:"f_crt"`
	PublicKey string    `json:"f_publicKey"`
	RootCrt   string    `json:"f_root_crt"`
	Created   time.Time `json:"f_created_at"`
}

func NewCert() *Cert {
	return &Cert{}
}

func (c *Cert) Insert(tx *sql.Tx) error {
	c.Created = time.Now()
	_, err := tx.Exec("insert into t_cert(f_chain_id, f_node_id, f_crt, f_publicKey, f_root_crt, f_created_at) values(?, ?, ?, ?, ?, ?)",
		c.ChainID, c.NodeID, c.Crt, c.PublicKey, c.RootCrt, c.Created)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cert) UpdateCert(tx *sql.Tx) error {
	res, err := tx.Exec("update t_cert set f_publicKey=?, f_crt=?, f_created_at=? where f_chain_id=? and f_node_id=?", c.PublicKey, c.Crt, c.Created, c.ChainID, c.NodeID)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (c *Cert) Delete(tx *sql.Tx, condition string) error {
	return nil
}

func (c *Cert) DeleteAll(tx *sql.Tx) error {
	if _, err := tx.Exec("delete from t_cert"); err != nil {
		return err
	}
	return nil
}
