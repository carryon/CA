package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/bocheninc/L0/components/crypto"
	"github.com/bocheninc/L0/components/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DeployServer string `yaml:"deploy_server"`
	BaseDir      string `yaml:"base_dir"`
	RootKey      string `yaml:"root_key"`
	RootCrt      string `yaml:"root_crt"`
	CACertDir    string `yaml:"ca_cert_dir"`
	AgentCertDir string `yaml:"agent_cert_dir"`
}

type ServerRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	Id     uint             `json:"id"`
}

type Params []struct {
	NodeID  string `json:"node_id"`
	AgentID string `json:"agent_id"`
}

type NodeCert struct {
	ID      uint `json:"id"`
	Results struct {
		NodeID   string `json:"node_id"`
		AgentKey string `json:"agent_key"`
		AgentCrt string `json:"agent_crt"`
		RootCrt  string `json:"root_crt"`
	} `json:"results"`
	Error interface{} `json:"error"`
}

type ConfigTimestamp struct {
	Result []struct {
		NodeID     string `json:"node_id"`
		UpdateTime int    `json:"update_time"`
	} `json:"result"`
	Error interface{} `json:"error"`
	ID    int         `json:"id"`
}

var conf = Config{}

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 5,
	},
}

func pathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func initConfigFile() {
	if !pathExist(conf.CACertDir) {
		err := os.MkdirAll(conf.CACertDir, os.ModePerm)
		if err != nil {
			log.Errorln("Create ca cert dir error:", err)
			return
		}
	}
	if *rootPtr {
		info := crypto.CertInformation{IsCA: true, CrtName: conf.RootCrt, KeyName: conf.RootKey}
		err := crypto.CreateCRT(nil, nil, info)
		if err != nil {
			log.Errorln("create root cert error:", err)
			return
		}
	} else {
		if !pathExist(conf.RootCrt) || !pathExist(conf.RootKey) {
			info := crypto.CertInformation{IsCA: true, CrtName: conf.RootCrt, KeyName: conf.RootKey}
			err := crypto.CreateCRT(nil, nil, info)
			if err != nil {
				log.Errorln("create root cert error:", err)
				return
			}
		}
	}
}

var rootPtr = flag.Bool("RegenCert", false, "regenerate cert file or use old cert file")

func init() {
	flag.Parse()

	logFile, err := os.Create("ca_server.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)

	//resolve config.yaml
	file, err := ioutil.ReadFile("config/ca_server.yaml")
	if err != nil {
		log.Errorln("read ca_server config error:", err)
		return
	}

	err = yaml.Unmarshal([]byte(file), &conf)
	if err != nil {
		log.Errorln("unmarshal ca_server config error:", err)
		return
	}

	initConfigFile()
}

func createCertFile(baseInfo crypto.CertInformation) (crtPath string, keyPath string, err error) {
	certDir := conf.AgentCertDir + baseInfo.Locality[0] + "/"
	if !pathExist(certDir) {
		err = os.Mkdir(certDir, os.ModePerm)
		if err != nil {
			log.Errorln("Create cert dir error:", err)
			return "", "", err
		}
	}
	crtName := baseInfo.Locality[0] + ".crt"
	keyName := baseInfo.Locality[0] + ".key"

	crtPath = certDir + crtName
	keyPath = certDir + keyName

	crt, pri, err := crypto.Parse(conf.RootCrt, conf.RootKey)
	if err != nil {
		log.Errorln("parse root crt and key error:", err)
		return "", "", err
	}
	crtInfo := baseInfo
	crtInfo.CrtName = crtPath
	crtInfo.KeyName = keyPath
	err = crypto.CreateCRT(crt, pri, crtInfo)
	if err != nil {
		log.Errorln("create crt error:", err)
		return "", "", err
	}

	return crtPath, keyPath, nil
}

func CaHandler(w http.ResponseWriter, r *http.Request) {
	reqObj := new(ServerRequest)
	err := json.NewDecoder(r.Body).Decode(reqObj)
	if err != nil {
		log.Errorln("decode request body error:", err)
		return
	}
	jsonByte, err := reqObj.Params.MarshalJSON()
	if err != nil {
		log.Errorln("marshal req params json error:", err)
		return
	}

	params := Params{}
	err = json.Unmarshal(jsonByte, &params)
	if err != nil {
		log.Errorln("unmarshal jsonByte error:", err)
		return
	}

	nodeID := params[0].NodeID
	agentID := params[0].AgentID

	// get all nodes info
	req, err := http.NewRequest("POST", conf.DeployServer, bytes.NewBufferString(
		`{"id":1,"method":"config-timestamp","params":["`+agentID+`"]}`,
	))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Errorln("post %s config-timestamp error:", conf.DeployServer, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("read config-timestamp response body error:", err)
		return
	}

	configTimestamp := ConfigTimestamp{}
	err = json.Unmarshal(body, &configTimestamp)
	if err != nil {
		log.Errorln("unmarshal config-timestamp response body error:", err)
		return
	}

	isValid := false
	for _, item := range configTimestamp.Result {
		if strings.EqualFold(item.NodeID, nodeID) {
			isValid = true
			break
		}
	}

	resObj := new(NodeCert)
	if !isValid {
		resObj.ID = reqObj.Id
		resObj.Error = "error: your node doesn't belong to your agent"

		ret, err := json.Marshal(resObj)
		if err != nil {
			log.Errorln("marshal resObj error:", err)
			return
		}
		w.Write(ret)
		return
	}

	baseInfo := crypto.CertInformation{Locality: []string{nodeID}}

	crtPath, keyPath, err := createCertFile(baseInfo)

	if err != nil {
		resObj.ID = reqObj.Id
		resObj.Error = "error: ca can't generate agent key and crt"

		ret, err := json.Marshal(resObj)
		if err != nil {
			log.Errorln("marshal resObj error:", err)
			return
		}
		w.Write(ret)
		return
	}

	agentCrt, err := ioutil.ReadFile(crtPath)
	if err != nil {
		log.Errorln("Read agentCrt file error, Error info:", err)
		return
	}

	agentKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Errorln("Read agentKey file error, Error info:", err)
		return
	}

	caCrt, err := ioutil.ReadFile(conf.RootCrt)
	if err != nil {
		log.Errorln("Read caCrt file error, Error info:", err)
		return
	}

	resObj.ID = reqObj.Id
	resObj.Error = nil
	resObj.Results.NodeID = nodeID
	resObj.Results.AgentCrt = string(agentCrt)
	resObj.Results.AgentKey = string(agentKey)
	resObj.Results.RootCrt = string(caCrt)

	ret, err := json.Marshal(resObj)
	if err != nil {
		log.Errorln("marshal resObj error:", err)
		return
	}
	w.Write(ret)
}
