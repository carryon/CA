package main

import (
	"fmt"
	"time"

	"github.com/bocheninc/CA/Agent-Server/config"
	"github.com/bocheninc/CA/Agent-Server/manager"
	"github.com/bocheninc/CA/Agent-Server/request"
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

func (s *Server) start() {
	//todo load exec file

	for {
		select {
		case <-s.ticker.C:
			nodes, err := s.req.GetLcndConfig(config.Cfg.ID)
			if err != nil {
				fmt.Println("Get nodes config err: ", err)
			}
			if err := s.manager.UpdateNodes(nodes); err != nil {
				fmt.Println("update nodes err: ", err)
			}
			if err := s.manager.StartNodes(); err != nil {
				fmt.Println("start nodes err: ", err)
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

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
	}

	config.Cfg = cfg

	s := NewServer()

	s.start()

}
