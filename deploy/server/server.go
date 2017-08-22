package server

import (
	"github.com/bocheninc/CA/deploy/components/db"
	"github.com/bocheninc/CA/deploy/components/log"
	"github.com/bocheninc/CA/deploy/config"
)

type Server struct {
	config *config.Config
	msgNet *MsgNet
	router *Router
}

func NewServer(configPath string) *Server {
	var s = new(Server)

	config, err := config.NewConfig(configPath)
	if err != nil {
		log.Error("load config file error:", err)
	}

	s.config = config

	//init log
	s.initLog()

	log.Infof("load config %v, db config %v", config, config.DBConfig)

	//init db
	db := db.NewDB(s.config.DBConfig)

	s.router = NewRouter(NewList(db), config.Port)

	return s
}

func (s *Server) Start() {
	log.Info("deploy server start...")
	s.router.start()
	// go s.msgNet.Start()

}

func (s *Server) initLog() {
	log.New(s.config.LogFile)
	log.SetLevel(s.config.LogLevel)
}
