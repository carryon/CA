package node

import (
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	yaml "gopkg.in/yaml.v2"

	"github.com/bocheninc/CA/Agent-Server/config"
	"github.com/bocheninc/CA/Agent-Server/log"
	"github.com/bocheninc/CA/Agent-Server/types"
	"github.com/bocheninc/CA/Agent-Server/utils"
)

type NodeInfo struct {
	sync.RWMutex
	NodeID     string
	Config     *types.NodeConfig
	Cert       *types.NodeCert
	CertPath   *types.CertConfig
	ConfigPath string
	IsRuning   bool
	Version    string
	cmd        *exec.Cmd
}

func NewNodeInfo(ID, version string, config *types.NodeConfig, cert *types.NodeCert) *NodeInfo {
	return &NodeInfo{NodeID: ID,
		Version:  version,
		Config:   config,
		Cert:     cert,
		CertPath: new(types.CertConfig),
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
		//start clear
		utils.RemoveFile(filepath.Join(config.Cfg.LcndDir, n.NodeID))

		if err := n.writeCert(); err != nil {
			return err
		}

		if err := n.writeConfig(); err != nil {
			return err
		}

		go n.startProcess(filepath.Join(config.Cfg.ExecDir, "lcnd"))

		n.IsRuning = true
	}
	return nil
}

func (n *NodeInfo) Stop() error {
	log.Infof("stop lcnd %s ,pid: %d and remove file", n.NodeID, n.cmd.Process.Pid)

	if err := n.cmd.Process.Kill(); err != nil {
		return err
	}

	return utils.RemoveFile(filepath.Join(config.Cfg.LcndDir, n.NodeID))
}

func (n *NodeInfo) writeConfig() error {

	nodeDir, err := utils.OpenDir(filepath.Join(config.Cfg.LcndDir, n.NodeID))
	if err != nil {
		return err
	}

	n.ConfigPath = filepath.Join(nodeDir, n.NodeID+".yaml")
	caConfig := new(types.CAConfig)
	if n.Cert == nil {
		caConfig.CA.Enabled = false
	} else {
		caConfig.CA.Enabled = true
	}

	caCnt, err := yaml.Marshal(caConfig)
	if err != nil {
		return err
	}

	n.Config.Vm.JsVMExeFilePath = filepath.Join(config.Cfg.ExecDir, "jsvm")
	n.Config.Vm.LuaVMExeFilePath = filepath.Join(config.Cfg.ExecDir, "luavm")

	nodeConfig, err := yaml.Marshal(n.Config)
	if err != nil {
		return err
	}

	cert, err := yaml.Marshal(n.CertPath)
	if err != nil {
		return err
	}

	f, err := utils.OpenFile(n.ConfigPath)
	if err != nil {
		return err
	}
	defer f.Close()

	//todo
	if _, err := f.Write(nodeConfig); err != nil {
		return err
	}

	if _, err := f.Write(cert); err != nil {
		return err
	}

	if _, err := f.Write(caCnt); err != nil {
		return err
	}

	return nil
}

func (n *NodeInfo) writeCert() error {
	nodeDir, err := utils.OpenDir(filepath.Join(config.Cfg.LcndDir, n.NodeID))
	if err != nil {
		return err
	}

	n.CertPath.Cert.CaPath = filepath.Join(nodeDir, "ca.crt")
	f, err := utils.OpenFile(n.CertPath.Cert.CaPath)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(n.Cert.Results.RootCrt); err != nil {
		return err
	}

	f.Close()

	n.CertPath.Cert.KeyPath = filepath.Join(nodeDir, n.NodeID+".key")
	f, err = utils.OpenFile(n.CertPath.Cert.KeyPath)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(n.Cert.Results.AgentKey); err != nil {
		return err
	}

	f.Close()

	n.CertPath.Cert.CrtPath = filepath.Join(nodeDir, n.NodeID+".crt")
	f, err = utils.OpenFile(n.CertPath.Cert.CrtPath)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(n.Cert.Results.AgentCrt); err != nil {
		return err
	}

	f.Close()

	return nil
}

func (n *NodeInfo) startProcess(execFilePath string) {
	log.Infof("start lcnd pocess: %s,config: %s", execFilePath, n.ConfigPath)
	n.cmd = exec.Command(execFilePath, "--config", n.ConfigPath)
	if err := n.cmd.Run(); err != nil {
		if n.cmd.ProcessState.Sys().(syscall.WaitStatus).Signal().String() == "killed" {
			log.Infof("lcnd %s was killed", n.NodeID)
			return
		}
		log.Errorf("start lcnd %s,err:%v", n.NodeID, err)
	}
}
