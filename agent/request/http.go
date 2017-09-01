package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bocheninc/CA/agent/config"
	"github.com/bocheninc/CA/agent/log"
	"github.com/bocheninc/CA/agent/manager/msgnet"
	"github.com/bocheninc/CA/agent/manager/node"
	"github.com/bocheninc/CA/agent/types"
)

type Resp struct {
	Result []string    `json:"result"`
	Error  interface{} `json:"error"`
	ID     uint        `json:"id"`
}

type Req struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     uint     `json:"id"`
}

type Request struct {
	client *http.Client
}

func NewRequest() *Request {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 5,
		},
	}

	return &Request{
		client: client,
	}
}

func (r *Request) getMsgnetVersion(id string) (string, error) {
	req := &Req{
		ID:     1,
		Method: "msgnet-version",
		Params: []string{id},
	}

	return r.getVersion(req, config.Cfg.DeployServer)
}

func (r *Request) getLcndVersion(id string) (string, error) {
	req := &Req{
		ID:     1,
		Method: "lcnd-version",
		Params: []string{id},
	}

	return r.getVersion(req, config.Cfg.DeployServer)
}

func (r *Request) GetLcndConfig(id string) ([]*node.NodeInfo, error) {
	var nodes []*node.NodeInfo

	req := &Req{
		ID:     1,
		Method: "nodes-config",
		Params: []string{id},
	}

	request, err := json.Marshal(req)
	if err != nil {
		log.Error("req err", err, string(request))
		return nil, err
	}

	data, err := r.request(config.Cfg.DeployServer, request)
	if err != nil {
		log.Error("data err", err, string(data))
		return nil, err
	}

	serverResponse := new(Resp)
	err = json.Unmarshal(data, &serverResponse)
	if err != nil {
		log.Error("serverResponse err", err, "----->", string(data))
		return nil, err
	}

	if serverResponse.Error != nil {
		return nil, fmt.Errorf("get nodes config resp: %s", serverResponse.Error.(string))
	}

	version, err := r.getLcndVersion(id)
	if err != nil {
		log.Error("get lcnf version ", err)
		return nil, err
	}

	for _, v := range serverResponse.Result {
		nodeConfig := new(types.NodeConfig)
		if err := json.Unmarshal([]byte(v), nodeConfig); err != nil {
			log.Error("node", err, v)
			return nil, err
		}
		cert, err := r.getCrt(id, nodeConfig.NodeID)
		if err != nil {
			log.Error("get crt err", err, *cert)
			return nil, err
		}
		node := node.NewNodeInfo(nodeConfig.NodeID, version, nodeConfig, cert)
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (r *Request) GetMsgnetConfig(id string) ([]*msgnet.MsgnetInfo, error) {
	var msgnets []*msgnet.MsgnetInfo

	req := &Req{
		ID:     1,
		Method: "msgnet-config",
		Params: []string{id},
	}

	request, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	data, err := r.request(config.Cfg.DeployServer, request)
	if err != nil {
		return nil, err
	}

	serverResponse := new(Resp)
	err = json.Unmarshal(data, &serverResponse)
	if err != nil {
		return nil, err
	}

	if serverResponse.Error != nil {
		return nil, errors.New(serverResponse.Error.(string))
	}

	version, err := r.getMsgnetVersion(id)
	if err != nil {
		return nil, err
	}

	for _, v := range serverResponse.Result {
		msgnetConfig := new(types.MsgnetConfig)
		if err := json.Unmarshal([]byte(v), msgnetConfig); err != nil {
			return nil, err
		}
		msgnetInfo := msgnet.NewMsgnetInfo(string(msgnetConfig.Router.ID), version, msgnetConfig)
		msgnets = append(msgnets, msgnetInfo)
	}

	return msgnets, nil

}

func (r *Request) getCrt(agentID, nodeID string) (*types.NodeCert, error) {
	req := &Req{
		ID:     1,
		Method: "node-cert",
		Params: []string{agentID, nodeID},
	}

	request, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	data, err := r.request(config.Cfg.CaServer, request)
	if err != nil {
		return nil, err
	}

	nodeCert := new(types.NodeCert)
	if err := json.Unmarshal(data, nodeCert); err != nil {
		return nil, err
	}

	if nodeCert.Error != nil {
		return nil, fmt.Errorf("get cert :%s", nodeCert.Error.(string))
	}

	return nodeCert, nil
}
func (r *Request) getVersion(req *Req, url string) (string, error) {

	request, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	data, err := r.request(url, request)
	if err != nil {
		return "", err
	}

	version := new(types.Version)
	err = json.Unmarshal(data, version)
	if err != nil {
		return "", err
	}

	if version.Error != nil {
		return "", errors.New(version.Error.(string))
	}
	return version.Result.Version, nil
}

func (r *Request) request(address string, request []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
