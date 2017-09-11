package sever

import (
	"time"

	"github.com/bocheninc/CA/agent/log"
	"github.com/bocheninc/CA/agent/manager"
)

type Server struct {
	ticker  *time.Ticker
	manager *manager.Manager
}

func NewServer() *Server {
	return &Server{
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
			if err := s.manager.UpdateNodes(); err != nil {
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
