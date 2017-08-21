package node

import (
	"path/filepath"
	"sync"

	yaml "gopkg.in/yaml.v2"

	"github.com/bocheninc/CA/Agent-Server/config"
	"github.com/bocheninc/CA/Agent-Server/types"
	"github.com/bocheninc/CA/Agent-Server/utils"
)

type NodeInfo struct {
	sync.RWMutex
	NodeID   string
	Config   *types.NodeConfig
	Cert     *types.NodeCert
	IsRuning bool
	Version  string
	//todo ca
}

func NewNodeInfo(ID, version string, config *types.NodeConfig, cert *types.NodeCert) *NodeInfo {
	return &NodeInfo{NodeID: ID,
		Version:  version,
		Config:   config,
		Cert:     cert,
		IsRuning: false}

}

func (n *NodeInfo) CheckVersion(version string) bool {
	if n.Version == version {
		return true
	}
	return false
}

func (n *NodeInfo) Start() error {

	if !n.IsRuning {
		configFilePath, err := n.writeConfig()
		if err != nil {
			return err
		}

		if err := n.writeCert(); err != nil {
			return err
		}

		if err := utils.StartProcess(filepath.Join(config.Cfg.ExecDir, "lcnd"), configFilePath); err != nil {
			return err
		}
		n.IsRuning = true
	}
	return nil
}

func (n *NodeInfo) Stop() error {
	return utils.StopProcess(n.NodeID)
}

func (n *NodeInfo) writeConfig() (string, error) {

	nodeDir, err := utils.OpenDir(filepath.Join(config.Cfg.LcndDir, n.NodeID))
	if err != nil {
		return "", err
	}

	configFilePath := filepath.Join(nodeDir, n.NodeID+".yaml")
	caConfig := new(types.CAConfig)
	if n.Cert == nil {
		caConfig.CA.Enabled = false
	} else {
		caConfig.CA.Enabled = true
	}

	caCnt, err := yaml.Marshal(caConfig)
	if err != nil {
		return "", err
	}

	nodeConfig, err := yaml.Marshal(n.Config)
	if err != nil {
		return "", err
	}

	cert, err := yaml.Marshal(n.Cert)
	if err != nil {
		return "", err
	}

	f, err := utils.OpenFile(configFilePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	//todo
	if _, err := f.Write(nodeConfig); err != nil {
		return "", err
	}

	if _, err := f.Write(cert); err != nil {
		return "", err
	}

	if _, err := f.Write(caCnt); err != nil {
		return "", err
	}

	return configFilePath, nil
}

func (n *NodeInfo) writeCert() error {
	nodeDir, err := utils.OpenDir(filepath.Join(config.Cfg.LcndDir, n.NodeID))
	if err != nil {
		return err
	}

	f, err := utils.OpenFile(filepath.Join(nodeDir, "ca.crt"))
	if err != nil {
		return err
	}

	if _, err := f.WriteString(n.Cert.Results.RootCrt); err != nil {
		return err
	}

	f.Close()

	f, err = utils.OpenFile(filepath.Join(nodeDir, n.NodeID+".key"))
	if err != nil {
		return err
	}

	if _, err := f.WriteString(n.Cert.Results.AgentKey); err != nil {
		return err
	}

	f.Close()

	f, err = utils.OpenFile(filepath.Join(nodeDir, n.NodeID+".crt"))
	if err != nil {
		return err
	}

	if _, err := f.WriteString(n.Cert.Results.AgentCrt); err != nil {
		return err
	}

	f.Close()

	return nil
}
