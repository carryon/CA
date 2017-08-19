package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"sync"
	// mysql driver
	"github.com/bocheninc/CA/deploy/components/log"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db     *sql.DB
	once   sync.Once
	config *MysqldbConfig
)

func open() *sql.DB {
	//open db
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8&loc=%s&parseTime=true",
		config.User, config.PWD, config.Host, config.Port, url.QueryEscape(config.Zone))
	sqldb, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Errorf("open mysql error connStr=%s err=%v", connStr, err)
	}

	// if not exists create database
	sqldb.SetMaxOpenConns(2000)
	sqldb.SetMaxIdleConns(1000)
	if err := sqldb.Ping(); err != nil {
		log.Error(err)
	}
	if _, err := sqldb.Exec(fmt.Sprintf("create database if not exists %s;", config.Name)); err != nil {
		log.Error(err)
	}
	if _, err := sqldb.Exec(fmt.Sprintf("use %s;", config.Name)); err != nil {
		log.Error(err)
	}

	// if not exists create tables
	//create agent table
	createAgentTableSQL := `CREATE TABLE IF NOT EXISTS t_agent (
	f_id BIGINT(20) UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	f_agent_id VARCHAR(255) NOT NULL COMMENT '终端ID',
	f_addr VARCHAR(255) NOT NULL COMMENT '地址',
	f_created_at TIMESTAMP DEFAULT NOW() COMMENT '创建时间',
	UNIQUE KEY uniq_agent_id (f_agent_id)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='终端表';`
	if _, err := sqldb.Exec(createAgentTableSQL); err != nil {
		log.Error("create Agent Table Sql error:", err)
	}

	//create node table
	CreateNodeTableSQL := `CREATE TABLE IF NOT EXISTS t_node (
	f_id BIGINT(20) UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	f_version INT(9) UNSIGNED DEFAULT 1 COMMENT '版本',
	f_chain_id VARCHAR(255) NOT NULL COMMENT '链ID',
	f_node_id VARCHAR(255) NOT NULL COMMENT '节点ID',
	f_agent_id VARCHAR(255) NOT NULL COMMENT '终端ID',
	f_config_file TEXT(10240) NOT NULL COMMENT '配置文件',
	f_config VARCHAR(255) NOT NULL COMMENT '区块配置',
	f_status VARCHAR(255) NOT NULL COMMENT '状态',
	f_height INT(9) NOT NULL COMMENT '高度',
	f_addr VARCHAR(255) NOT NULL COMMENT '地址',
	f_created_at DATETIME COMMENT '创建时间',
	f_updated_at DATETIME COMMENT '更新时间',
	UNIQUE KEY uniq_chain_id_node_id (f_chain_id,f_node_id)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='节点表';`
	if _, err := sqldb.Exec(CreateNodeTableSQL); err != nil {
		log.Error("create Node Table Sql error:", err)
	}
	return sqldb
}

//NewDB return a MySQL object
func NewDB(c *MysqldbConfig) *sql.DB {
	once.Do(func() {
		config = c
		db = open()
	})
	return db
}
