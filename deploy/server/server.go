package server

import (
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/bocheninc/CA/deploy/components/db"
	"github.com/bocheninc/CA/deploy/components/log"
	"github.com/bocheninc/CA/deploy/config"
)

var defaultAddr = "deploy:server"

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

	s.msgNet = NewMsgnet(defaultAddr, s.config.RouterAddresses, s.handleMsgnetMessage, s.config.LogDir)

	return s
}

func (s *Server) Start() {
	log.Info("deploy server start...")
	go s.msgNet.Start()
	go s.changeNodeConfigLoop()
	s.router.start()

}

func (s *Server) changeNodeConfigLoop() {
	for {
		select {
		case msg := <-s.router.list.msgChan:
			s.msgNet.SendMsgnetMessage(defaultAddr, msg.ChainID+":", msg.Message)
		}
	}
}

func (s *Server) initLog() {
	log.New(s.config.LogFile)
	log.SetLevel(s.config.LogLevel)
}

func (s *Server) handleMsgnetMessage(src, dst string, payload, signature []byte) error {

	msg := Message{}
	msg.Deserialize(payload)
	log.Debugf("recv msg-net message src: %s, dst: %s,type %d ,payload: %s .", src, dst, msg.Cmd, payload)

	switch msg.Cmd {
	case ChainNodeStatusMsg:
		chainID, peerID := parseID(src)
		node, err := s.router.list.GetNodeByNodeID(chainID, peerID)
		if err != nil {
			log.Error("get node info err: ", err)
			return err
		}

		type Status struct {
			Height int
			Tps    int
		}

		status := new(Status)

		json.Unmarshal(msg.Payload, status)

		node.Height = status.Height
		node.Status = strconv.Itoa(status.Tps)
		node.Updated = time.Now()

		tx, _ := s.router.list.Db.Begin()
		if err := node.UpdateHeight(tx); err != nil {
			log.Errorf(" chainID: %s, nodeID: %s update height: %d ,err: %v", node.ChainID, node.NodeID, node.Height, err)
			tx.Rollback()
			return err
		}
		tx.Commit()

		log.Debugln("recv node report status ", chainID, peerID, node.Height, node.Status)
	default:
		log.Warn("recv not know msgnet.type: ", msg.Cmd)
	}

	return nil
}

// parseID returns chainID and PeerID
func parseID(peerAddress string) (string, string) {
	id := strings.Split(peerAddress, ":")
	// chainID, _ := hex.DecodeString(id[0])
	// if len(id) == 2 {
	peerid, _ := hex.DecodeString(id[1])
	// 	return chainID, p2p.PeerID(peerid)
	// }
	//return chainID, nil
	return id[0], string(peerid)
}
