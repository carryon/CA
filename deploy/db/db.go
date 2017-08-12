package db

import (
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

type DeployDB struct {
	DB   *leveldb.DB
	Path string
}

var DefaultPath string = "/tmp/leveldb/"

func init() {
	if _, err := os.Stat(DefaultPath); err != nil {
		os.Mkdir(DefaultPath, 0755)
	}
}
func NewDB(path string) (*DeployDB, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	dd := new(DeployDB)
	dd.DB = db
	dd.Path = path

	return dd, nil
}

func (self *DeployDB) Put(key, value []byte) error {
	return self.DB.Put(key, value, nil)
}

func (self *DeployDB) Get(key []byte) ([]byte, error) {
	return self.DB.Get(key, nil)
}

func (self *DeployDB) Delete(key []byte) error {
	return self.DB.Delete(key, nil)
}
