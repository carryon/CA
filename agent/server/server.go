package sever

import (
	"time"

	"github.com/bocheninc/CA/agent/config"
	"github.com/bocheninc/CA/agent/log"
	"github.com/bocheninc/CA/agent/manager"
	"github.com/bocheninc/CA/agent/request"
)

type Server struct {
	req     *request.Request
	ticker  *time.Ticker
	manager *manager.Manager
}

func NewServer() *Server {
	return &Server{
		req:     request.NewRequest(),
		ticker:  time.NewTicker(10 * time.Second),
		manager: manager.NewManager(),
	}
}

func (s *Server) Start() {
	log.Info("agent server start ...")
	//todo load exec file

	for {
		select {
		case <-s.ticker.C:
			nodes, err := s.req.GetLcndConfig(config.Cfg.ID)
			if err != nil {
				log.Error("Get nodes config err: ", err)
			}
			if err := s.manager.UpdateNodes(nodes); err != nil {
				log.Error("update nodes err: ", err)
			}
			if err := s.manager.StartNodes(); err != nil {
				log.Error("start nodes err: ", err)
			}

			//todo msgnet
			// msgnets,err:=s.req.GetMsgnetConfig(config.Cfg.ID)
			// if err!=nil{
			// 	fmt.Println("Get msgnet config err: ",err)
			// }
			// if err:=s.manager.UpdateNodes(nodes);err!=nil{
			// 	fmt.Println("update msgnets err: ",err)
			// }
			// if err:=s.manager.StartMsgnets();err!=nil{
			// 	fmt.Println("start msgnet err: ",err)
			// }

		}
	}

}
