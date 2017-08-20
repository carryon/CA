package node

import (
	"sync"
)

type NodeInfo struct{
	sync.RWMutex
	NodeID string
	Config []byte
	IsRuning bool
	Version string
	//todo ca
}


func NewNodeInfo(ID ,version string,config []byte)*NodeInfo{
	return &NodeInfo{NodeID:ID,
		Version:version,
		Config:config,
		IsRuning:false}
		
}

func (n*NodeInfo)CheckVersion(version string)bool{
	if n.Version == version{
		return true
	}
	return false
}

func (n*NodeInfo)Start()error{
	if !n.IsRuning {
		

		n.IsRuning = true	
	}
	return nil
}

func (n*NodeInfo)Stop()error{
	return nil
}