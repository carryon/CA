package node

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	yaml "gopkg.in/yaml.v2"

	"github.com/bocheninc/CA/agent/config"
	"github.com/bocheninc/CA/agent/log"
	"github.com/bocheninc/CA/agent/types"
	"github.com/bocheninc/CA/agent/utils"
)

type NodeInfo struct {
	sync.RWMutex
	NodeID     string
	Config     *types.NodeConfig
	Cert       *types.NodeCert
	CaConfig   *types.CAConfig
	ConfigPath string
	IsRuning   bool
	Version    string
	cmd        *exec.Cmd
	nodeDir    string
}

func NewNodeInfo(ID, version string, nodeConfig *types.NodeConfig, cert *types.NodeCert) *NodeInfo {
	return &NodeInfo{NodeID: ID,
		Version:  version,
		Config:   nodeConfig,
		Cert:     cert,
		CaConfig: new(types.CAConfig),
		IsRuning: false,
	}
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

		nodeDir, err := utils.OpenDir(filepath.Join(config.Cfg.LcndDir, n.NodeID))
		if err != nil {
			log.Error(n.NodeID, " opendir error :", err)
			return err
		}
		n.nodeDir = nodeDir

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
		log.Error("kill ", n.NodeID, " error ", err)
	}

	return utils.RemoveFile(filepath.Join(config.Cfg.LcndDir, n.NodeID))
}

func (n *NodeInfo) writeCert() error {
	//root
	n.CaConfig.CA.Cert.CaPath = filepath.Join(n.nodeDir, "ca.crt")

	f, err := utils.OpenFile(n.CaConfig.CA.Cert.CaPath)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(n.Cert.RootCertificate); err != nil {
		return err
	}
	f.Close()

	n.CaConfig.CA.Cert.KeyPath = filepath.Join(n.nodeDir, n.NodeID+".key")
	key := x509.MarshalPKCS1PrivateKey(n.Cert.PrivateKey)
	if err := write(n.CaConfig.CA.Cert.KeyPath, "PRIVATE KEY", key); err != nil {
		return fmt.Errorf("write private key err: %s", err)

	}

	n.CaConfig.CA.Cert.CrtPath = filepath.Join(n.nodeDir, n.NodeID+".crt")
	f, err = utils.OpenFile(n.CaConfig.CA.Cert.CrtPath)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(n.Cert.Certificate); err != nil {
		return err
	}
	f.Close()

	return nil
}

func (n *NodeInfo) writeConfig() error {
	n.ConfigPath = filepath.Join(n.nodeDir, n.NodeID+".yaml")

	if n.Cert == nil {
		n.CaConfig.CA.Enabled = false
	} else {
		n.CaConfig.CA.Enabled = true
	}

	caCnt, err := yaml.Marshal(n.CaConfig)
	if err != nil {
		return err
	}

	n.Config.Blockchain.Datadir = filepath.Join(config.Cfg.BaseDir, n.NodeID)
	n.Config.Vm.JsVMExeFilePath = filepath.Join(config.Cfg.ExecDir, "jsvm")
	n.Config.Vm.LuaVMExeFilePath = filepath.Join(config.Cfg.ExecDir, "luavm")

	nodeConfig, err := yaml.Marshal(n.Config)
	if err != nil {
		return err
	}

	var b bytes.Buffer

	b.Write(nodeConfig)
	b.Write(caCnt)

	ioutil.WriteFile(n.ConfigPath, b.Bytes(), 0666)

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

func write(filename, Type string, data []byte) error {
	if utils.FileExist(filename) {
		if err := os.RemoveAll(filename); err != nil {
			return err
		}
	}
	File, err := os.Create(filename)
	defer File.Close()
	if err != nil {
		return err
	}
	return pem.Encode(File, &pem.Block{Bytes: data, Type: Type})
}
