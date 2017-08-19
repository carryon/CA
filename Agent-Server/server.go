package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bocheninc/L0/components/log"
	util "github.com/ghodss/yaml"
	yaml "gopkg.in/yaml.v2"
)

var lcnd_version = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

var l0NodeInfo = struct {
	sync.RWMutex
	m map[string]int
}{m: make(map[string]int)}

var msgnet_version = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

var msgnetNodeInfo = struct {
	sync.RWMutex
	m map[string]int
}{m: make(map[string]int)}

var conf = Config{}
var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 5,
	},
}

func initConfigFile() {
	if !pathExist(conf.AgentCertDir) {
		err := os.MkdirAll(conf.AgentCertDir, os.ModePerm)
		if err != nil {
			log.Errorln("Create agent cert dir error:", err)
			return
		}
	}

	if !pathExist(conf.MsgnetConfigDir) {
		err := os.MkdirAll(conf.MsgnetConfigDir, os.ModePerm)
		if err != nil {
			log.Errorln("Create msg-net config dir error:", err)
			return
		}
	}
}

func init() {
	logFile, err := os.Create("agent_server.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)

	data, err := ioutil.ReadFile("config/agent.yaml")
	if err != nil {
		log.Errorln("read agent config file error:", err)
		return
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Errorln("unmarshal error, error:", err)
		return
	}
	initConfigFile()
	msgnetInit()
	time.Sleep(time.Second * 5)
	l0Init()
}

func main() {
	for {
		checkMsgnet()
		checkL0()
		time.Sleep(time.Minute * 1)
	}
}

func checkL0() {
	checkL0Nodes()
	checkL0Version()
}

func checkMsgnet() {
	checkMsgnetNodes()
	checkMsgnetVersion()
}

func l0Init() {
	getLcndVersion()
	if !pathExist(conf.L0ExecFile) {
		downloadLcnd()
	}
	startL0NodesService(nil)
}

func msgnetInit() {
	getMsgnetVersion()
	if !pathExist(conf.MsgnetExecFile) {
		downloadMsgnet()
	}
	startMsgnetNodesService(nil)
}

func getMsgnetVersion() {
	req, err := http.NewRequest("POST", conf.DeployServer, bytes.NewBufferString(
		`{"id":1,"method":"msgnet-version","params":["`+conf.ID+`"]}`,
	))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("post %v msgnet-version error: %v\n", conf.DeployServer, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read msgnet-version response body error: %v\n", err)
		return
	}

	msgnetVersion := VersionAPI{}
	err = json.Unmarshal(body, &msgnetVersion)
	if err != nil {
		log.Errorf("unmarshal msgnet-version response body error: %v\n", err)
		return
	}

	apiErr := msgnetVersion.Error
	if apiErr != nil {
		log.Errorf("can't get msgnet version, error: %v\n", err)
		return
	}

	msgnet_version.Lock()
	msgnet_version.m["version"] = msgnetVersion.Result.Version
	msgnet_version.Unlock()
}

func downloadMsgnet() {
	downloadCmd := fmt.Sprintf("wget -P %s %s", conf.BaseDir, conf.MsgnetUrl)
	cmd := exec.Command("/bin/sh", "-c", downloadCmd)
	err := cmd.Run()
	if err != nil {
		log.Errorln("execute downloadCmd:", downloadCmd, "error:", err)
		return
	}

	msgnetExecFile := conf.MsgnetExecFile
	chmodCmd := fmt.Sprintf("chmod a+x %s", msgnetExecFile)
	cmd = exec.Command("/bin/sh", "-c", chmodCmd)
	err = cmd.Run()
	if err != nil {
		log.Errorln("execute chmodCmd:", chmodCmd, "error:", err)
		return
	}
}

func startMsgnetNodesService(filterNodes []string) {
	req, err := http.NewRequest("POST", conf.DeployServer, bytes.NewBufferString(
		`{"id":1,"method":"msgnet-config","params":["`+conf.ID+`"]}`,
	))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("post %s msgnet-config error: %v\n", conf.DeployServer, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Errorf("read msgnet-config response body error: %v\n", err)
		return
	}

	serverResponse := ServerResponse{}
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		log.Errorf("unmarshal msgnet-config response body error: %v\n", err)
		return
	}

	results := serverResponse.Result
	for _, result := range results {
		msgnetConfig := MsgnetConfig{}
		err = json.Unmarshal([]byte(result), &msgnetConfig)
		if err != nil {
			log.Errorf("unmarshal result error: %v\n", err)
			return
		}

		isContinue := false
		for _, filterNode := range filterNodes {
			if strings.EqualFold(filterNode, msgnetConfig.NodeID) {
				isContinue = true
				break
			}
		}

		if isContinue {
			continue
		}

		// generate msg-net config file
		nodeID := msgnetConfig.NodeID
		configFileName := conf.MsgnetConfigDir + nodeID + ".yaml"
		confContent, err := util.JSONToYAML([]byte(result))
		if err != nil {
			log.Errorf("convert msg-net config json to yaml error: %v\n", err)
			return
		}
		writeMsgnetConfigFile(configFileName, confContent)

		// start msg-net
		startMsgnetService(configFileName)

		msgnetNodeInfo.Lock()
		msgnetNodeInfo.m[msgnetConfig.NodeID] = msgnetConfig.UpdateTime
		msgnetNodeInfo.Unlock()
	}
}

func getLcndVersion() {
	req, err := http.NewRequest("POST", conf.DeployServer, bytes.NewBufferString(
		`{"id":1,"method":"lcnd-version","params":["`+conf.ID+`"]}`,
	))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("post %v lcnd-version error: %v\n", conf.DeployServer, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read lcnd-version response body error: %v\n", err)
		return
	}

	lcndVersion := VersionAPI{}
	err = json.Unmarshal(body, &lcndVersion)
	if err != nil {
		log.Errorf("unmarshal lcnd-version response body error: %v\n", err)
		return
	}

	apiErr := lcndVersion.Error
	if apiErr != nil {
		log.Errorf("can't get lcnd version, error: %v\n", err)
		return
	}

	lcnd_version.Lock()
	lcnd_version.m["version"] = lcndVersion.Result.Version
	lcnd_version.Unlock()
}

func downloadLcnd() {
	downloadCmd := fmt.Sprintf("wget -P %s %s", conf.BaseDir, conf.LcndUrl)
	cmd := exec.Command("/bin/sh", "-c", downloadCmd)
	err := cmd.Run()
	if err != nil {
		log.Errorln("execute downloadCmd:", downloadCmd, "error:", err)
		return
	}

	lcndExecFile := conf.L0ExecFile
	chmodCmd := fmt.Sprintf("chmod a+x %s", lcndExecFile)
	cmd = exec.Command("/bin/sh", "-c", chmodCmd)
	err = cmd.Run()
	if err != nil {
		log.Errorln("execute chmodCmd:", chmodCmd, "error:", err)
		return
	}
}

func startL0NodesService(filterNodes []string) {
	req, err := http.NewRequest("POST", conf.DeployServer, bytes.NewBufferString(
		`{"id":1,"method":"nodes-config","params":["`+conf.ID+`"]}`,
	))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("post %s node-config error: %v\n", conf.DeployServer, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Errorf("read node-config response body error: %v\n", err)
		return
	}

	serverResponse := ServerResponse{}
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		log.Errorf("unmarshal node-config response body error: %v\n", err)
		return
	}

	results := serverResponse.Result
	for _, result := range results {
		nodeConfig := NodeConfig{}
		err = json.Unmarshal([]byte(result), &nodeConfig)
		if err != nil {
			log.Errorf("unmarshal result error: %v\n", err)
			return
		}

		// filter nodes
		isContinue := false
		for _, filterNode := range filterNodes {
			if strings.EqualFold(filterNode, nodeConfig.NodeID) {
				isContinue = true
				break
			}
		}

		if isContinue {
			continue
		}

		// get nodes cert
		req, err := http.NewRequest("POST", conf.CaServer, bytes.NewBufferString(
			`{"id":1,"method":"node-cert","params":[{"agent_id": "`+conf.ID+`", "node_id": "`+nodeConfig.NodeID+`"}]}`,
		))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("post %s node-cert error: %v\n", conf.CaServer, err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Errorf("read node-cert response body error: %v\n", err)
			return
		}

		nodeCert := NodeCert{}
		err = json.Unmarshal(body, &nodeCert)
		if err != nil {
			log.Errorf("unmarshal node-cert response body error: %v\n", err)
			return
		}

		apiErr := nodeCert.Error
		if apiErr != nil {
			log.Errorf("can't get node cert, error: %v\n", apiErr)
			return
		}

		// write cert content
		keyContent := nodeCert.Results.AgentKey
		crtContent := nodeCert.Results.AgentCrt
		caContent := nodeCert.Results.RootCrt

		crtName := nodeConfig.NodeID + ".crt"
		keyName := nodeConfig.NodeID + ".key"
		l0ConfigName := nodeConfig.NodeID + ".yaml"

		certDir := conf.AgentCertDir + nodeConfig.NodeID

		crtPath := certDir + "/" + crtName
		keyPath := certDir + "/" + keyName
		caPath := certDir + "/" + "ca.crt"
		l0ConfigPath := certDir + "/" + l0ConfigName

		confContent, err := util.JSONToYAML([]byte(result))
		if err != nil {
			log.Errorf("convert l0 config json to yaml error: %v\n", err)
			return
		}

		writeCertConfigFile(certDir, crtPath, keyPath, caPath, l0ConfigPath, crtContent, keyContent, caContent, confContent)

		// modify l0 config file
		certConfig := CertConfig{}
		certConfig.Cert.KeyPath = keyPath
		certConfig.Cert.CrtPath = crtPath
		certConfig.Cert.CaPath = caPath

		certContent, err := yaml.Marshal(certConfig)
		if err != nil {
			log.Errorf("marshal certConfig error: %v\n", err)
			return
		}

		caConfig := CAConfig{}
		if len(caContent) == 0 {
			caConfig.CA.Enabled = false
		} else {
			caConfig.CA.Enabled = true
		}

		caCnt, err := yaml.Marshal(caConfig)
		if err != nil {
			log.Errorf("marshal caConfig error: %v\n", err)
			return
		}

		f, err := os.OpenFile(l0ConfigPath, os.O_APPEND|os.O_RDWR, os.ModeAppend)
		if err != nil {
			log.Errorf("open l0ConfigPath: %v error: %v\n", l0ConfigPath, err)
			return
		}
		defer f.Close()
		_, err = f.Write(certContent)
		if err != nil {
			log.Errorf("write certConfig error: %v\n", err)
			return
		}
		_, err = f.Write(caCnt)
		if err != nil {
			log.Errorf("write caConfig error: %v\n", err)
			return
		}

		startL0Service(l0ConfigPath)

		l0NodeInfo.Lock()
		l0NodeInfo.m[nodeConfig.NodeID] = nodeConfig.UpdateTime
		l0NodeInfo.Unlock()
	}
}

func checkL0Version() {
	// get l0 version
	req, err := http.NewRequest("POST", conf.DeployServer, bytes.NewBufferString(
		`{"id":1,"method":"lcnd-version","params":["`+conf.ID+`"]}`,
	))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("request %v lcnd-version error: %v\n", conf.DeployServer, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Errorf("read lcnd-version response body error: %v\n", err)
		return
	}

	lcndVersion := VersionAPI{}
	err = json.Unmarshal(body, &lcndVersion)
	if err != nil {
		log.Errorf("unmarshal lcnd-version response body error: %v\n", err)
		return
	}

	apiErr := lcndVersion.Error
	if apiErr != nil {
		log.Errorf("can't get lcnd version, error: %v\n", apiErr)
		return
	}

	var nowVersion string
	lcnd_version.RLock()
	nowVersion = lcnd_version.m["version"]
	lcnd_version.RUnlock()

	if !strings.EqualFold(nowVersion, lcndVersion.Result.Version) {
		stopService("lcnd")
		removeLcndFile()
		getLcndVersion()
		downloadLcnd()
		configs := findConfigsPath()
		for _, config := range configs {
			startL0Service(config)
		}
	}
}

func checkMsgnetVersion() {
	req, err := http.NewRequest("POST", conf.DeployServer, bytes.NewBufferString(
		`{"id":1,"method":"msgnet-version","params":["`+conf.ID+`"]}`,
	))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("request %v msgnet-version error: %v\n", conf.DeployServer, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Errorf("read msgnet-version response body error: %v\n", err)
		return
	}

	msgnetVersion := VersionAPI{}
	err = json.Unmarshal(body, &msgnetVersion)
	if err != nil {
		log.Errorf("unmarshal msgnet-version response body error: %v\n", err)
		return
	}

	apiErr := msgnetVersion.Error
	if apiErr != nil {
		log.Errorf("can't get msg-net version, error: %v\n", apiErr)
		return
	}

	var nowVersion string
	msgnet_version.RLock()
	nowVersion = msgnet_version.m["version"]
	msgnet_version.RUnlock()

	if !strings.EqualFold(nowVersion, msgnetVersion.Result.Version) {
		stopService("msg-net")
		removeMsgnetFile()
		getMsgnetVersion()
		downloadMsgnet()
		configs := findMsgnetConfigsPath()
		for _, config := range configs {
			startMsgnetService(config)
		}
	}
}

func checkL0Nodes() {
	// get nodes timestamp
	req, err := http.NewRequest("POST", conf.DeployServer, bytes.NewBufferString(
		`{"id":1,"method":"config-timestamp","params":["`+conf.ID+`"]}`,
	))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("request %v config-timestamp error: %v\n", conf.DeployServer, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Errorf("read config-timestamp response body error: %v\n", err)
		return
	}

	configTimestamp := TimestampAPI{}
	err = json.Unmarshal(body, &configTimestamp)
	if err != nil {
		log.Errorf("unmarshal config-timestamp response body error: %v\n", err)
		return
	}

	apiErr := configTimestamp.Error
	result := configTimestamp.Result

	if apiErr != nil || result == nil {
		log.Errorln("can't get nodes info, error:", apiErr)
		stopService("lcnd")
		clearNodeInfo()
		removeDir(conf.AgentCertDir)
	}

	//find diff nodes and same nodes
	var diffNodes []string
	var sameNodes []string
	isSame := false
	l0NodeInfo.RLock()
	for lk, lv := range l0NodeInfo.m {
		isSame = false
		for _, cv := range result {
			if strings.EqualFold(lk, cv.NodeID) && lv == cv.UpdateTime {
				isSame = true
				break
			}
		}
		if isSame {
			sameNodes = append(sameNodes, lk)
		} else {
			diffNodes = append(diffNodes, lk)
		}
	}
	l0NodeInfo.RUnlock()

	// stop all diff nodes
	for _, v := range diffNodes {
		stopService(v)
	}
	removeStopNodeDir(diffNodes)
	l0NodeInfo.Lock()
	for _, k := range diffNodes {
		delete(l0NodeInfo.m, k)
	}
	l0NodeInfo.Unlock()

	// start config nodes
	startL0NodesService(sameNodes)
}

func checkMsgnetNodes() {
	// get nodes timestamp
	req, err := http.NewRequest("POST", conf.DeployServer, bytes.NewBufferString(
		`{"id":1,"method":"msgnet-timestamp","params":["`+conf.ID+`"]}`,
	))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("request %v msgnet-timestamp error: %v\n", conf.DeployServer, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Errorf("read msgnet-timestamp response body error: %v\n", err)
		return
	}

	msgnetTimestamp := TimestampAPI{}
	err = json.Unmarshal(body, &msgnetTimestamp)
	if err != nil {
		log.Errorf("unmarshal msgnet-timestamp response body error: %v\n", err)
		return
	}

	apiErr := msgnetTimestamp.Error
	result := msgnetTimestamp.Result

	if apiErr != nil || result == nil {
		log.Errorln("can't get msg-net nodes info, error:", apiErr)
		stopService("msg-net")
		clearMsgnetNodeInfo()
		removeDir(conf.MsgnetConfigDir)
	}

	// if err == nil && result == nil {
	// 	log.Errorln("no nodes")
	// 	stopService("lcnd")
	// 	clearNodeInfo()
	// 	removeDir(conf.AgentCertDir)
	// }

	//find diff nodes and same nodes
	var diffNodes []string
	var sameNodes []string
	isSame := false
	msgnetNodeInfo.RLock()
	for lk, lv := range msgnetNodeInfo.m {
		isSame = false
		for _, cv := range result {
			if strings.EqualFold(lk, cv.NodeID) && lv == cv.UpdateTime {
				isSame = true
				break
			}
		}
		if isSame {
			sameNodes = append(sameNodes, lk)
		} else {
			diffNodes = append(diffNodes, lk)
		}
	}
	msgnetNodeInfo.RUnlock()

	// stop all diff nodes
	for _, v := range diffNodes {
		stopService(v)
	}
	removeMegnetConfigFile(diffNodes)
	msgnetNodeInfo.Lock()
	for _, k := range diffNodes {
		delete(msgnetNodeInfo.m, k)
	}
	msgnetNodeInfo.Unlock()

	// start config nodes
	startMsgnetNodesService(sameNodes)
}

func startL0Service(l0Conf string) {
	l0ExecFile := conf.L0ExecFile
	startCmd := fmt.Sprintf("%s --config=%s &", l0ExecFile, l0Conf)

	cmd := exec.Command("/bin/sh", "-c", startCmd)
	err := cmd.Run()
	if err != nil {
		log.Errorln("execute startCmd:", startCmd, "error:", err)
		return
	}
}

func startMsgnetService(msgnetConf string) {
	msgnetExecFile := conf.MsgnetExecFile
	startCmd := fmt.Sprintf("%s --config=%s router &", msgnetExecFile, msgnetConf)

	cmd := exec.Command("/bin/sh", "-c", startCmd)
	err := cmd.Run()
	if err != nil {
		log.Errorln("execute startCmd:", startCmd, "error:", err)
		return
	}
}

func stopService(keyword string) {
	killCmd := fmt.Sprintf("pkill -f  %s", keyword)
	cmd := exec.Command("/bin/sh", "-c", killCmd)
	cmd.Run()
}

func removeDir(walkDir string) {
	var files []string
	err := filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return err
	})
	if err != nil {
		log.Errorln("walk error:", err)
		return
	}

	for _, file := range files {
		if strings.EqualFold(file, walkDir) {
			continue
		}
		os.RemoveAll(file)
	}
}

func clearNodeInfo() {
	l0NodeInfo.Lock()
	for k := range l0NodeInfo.m {
		delete(l0NodeInfo.m, k)
	}
	l0NodeInfo.Unlock()
}

func clearMsgnetNodeInfo() {
	msgnetNodeInfo.Lock()
	for k := range msgnetNodeInfo.m {
		delete(msgnetNodeInfo.m, k)
	}
	msgnetNodeInfo.Unlock()
}

func removeStopNodeDir(diffNodes []string) {
	var dirs []string
	for _, v := range diffNodes {
		dir := conf.AgentCertDir + v
		dirs = append(dirs, dir)
	}
	for _, dir := range dirs {
		os.RemoveAll(dir)
	}
}

func removeMegnetConfigFile(diffNodes []string) {
	dir := conf.MsgnetConfigDir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Errorln("read msg-net config dir error:", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			for _, diffNode := range diffNodes {
				if strings.Contains(file.Name(), diffNode) {
					err = os.Remove(dir + file.Name())
					if err != nil {
						log.Errorln("remove msg-net config file error:", err)
					}
				}
			}
		}
	}
}

func findConfigsPath() []string {
	var configArr []string
	l0NodeInfo.RLock()
	for node, _ := range l0NodeInfo.m {
		config := conf.AgentCertDir + node + "/" + node + ".yaml"
		configArr = append(configArr, config)
	}
	l0NodeInfo.RUnlock()
	return configArr
}

func findMsgnetConfigsPath() []string {
	var configArr []string
	msgnetNodeInfo.RLock()
	for node, _ := range msgnetNodeInfo.m {
		config := conf.MsgnetConfigDir + node + ".yaml"
		configArr = append(configArr, config)
	}
	msgnetNodeInfo.RUnlock()
	return configArr
}

func pathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func writeMsgnetConfigFile(configPath string, confContent []byte) {
	err := ioutil.WriteFile(configPath, []byte(confContent), 0666)
	if err != nil {
		log.Errorln("write msg-net config file:", configPath, "error:", err)
		return
	}
}

func writeCertConfigFile(certDir string, crtPath string, keyPath string, caPath string, l0ConfigPath string,
	crtContent string, keyContent string, caContent string, confContent []byte) {
	if pathExist(crtPath) {
		err := ioutil.WriteFile(crtPath, []byte(crtContent), 0666)
		if err != nil {
			log.Errorln("write crtPath:", crtPath, "error:", err)
			return
		}
	} else {
		err := os.Mkdir(certDir, os.ModePerm)
		if err != nil {
			log.Errorln("create certDir:", certDir, "error:", err)
			return
		}
		err = ioutil.WriteFile(crtPath, []byte(crtContent), 0666)
		if err != nil {
			log.Errorln("write crtPath:", crtPath, "error:", err)
			return
		}
	}

	err := ioutil.WriteFile(keyPath, []byte(keyContent), 0666)
	if err != nil {
		log.Errorln("write keyPath:", keyPath, "error:", err)
		return
	}

	err = ioutil.WriteFile(caPath, []byte(caContent), 0666)
	if err != nil {
		log.Errorln("write caPath:", caPath, "error:", err)
		return
	}

	err = ioutil.WriteFile(l0ConfigPath, []byte(confContent), 0666)
	if err != nil {
		log.Errorln("write l0ConfigPath:", l0ConfigPath, "error:", err)
		return
	}
}

func removeLcndFile() {
	lcnd := conf.L0ExecFile
	err := os.Remove(lcnd)
	if err != nil {
		log.Errorln("can't remove lcnd, error:", err)
		return
	}
}

func removeMsgnetFile() {
	msgnet := conf.MsgnetExecFile
	err := os.Remove(msgnet)
	if err != nil {
		log.Errorln("can't remove msg-net, error:", err)
		return
	}
}
