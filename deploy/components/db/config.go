package db

//MysqldbConfig db config
type MysqldbConfig struct {
	Name string
	User string
	PWD  string
	Host string
	Port string
	Zone string
}

// DefaultMysqldbConfig defines the default configuration of the mysqldb
func DefaultMysqldbConfig() *MysqldbConfig {
	return &MysqldbConfig{
		Name: "db_l0",
		User: "root",
		PWD:  "root",
		Host: "127.0.0.1",
		Port: "3306",
		Zone: "Asia/Shanghai",
	}
}
