package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bocheninc/CA/Agent-Server/config"
	"github.com/bocheninc/CA/Agent-Server/manager/msgnet"
	"github.com/bocheninc/CA/Agent-Server/manager/node"
	"github.com/bocheninc/CA/Agent-Server/types"
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

	return r.getVersion(req, config.Cfg.MsgnetURL)
}

func (r *Request) getLcndVersion(id string) (string, error) {
	req := &Req{
		ID:     1,
		Method: "lcnd-version",
		Params: []string{id},
	}

	return r.getVersion(req, config.Cfg.LcndURL)
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
		return nil, err
	}

	data, err := r.request(config.Cfg.LcndURL, request)
	if err != nil {
		return nil, err
	}

	serverResponse := new(Resp)
	err = json.Unmarshal(data, &serverResponse)
	if err != nil {
		return nil, err
	}

	version, err := r.getLcndVersion(id)
	if err != nil {
		return nil, err
	}

	for _, v := range serverResponse.Result {
		nodeConfig := new(types.NodeConfig)
		if err := json.Unmarshal([]byte(v), nodeConfig); err != nil {
			return nil, err
		}

		cert, err := r.getCrt(id, nodeConfig.NodeID)
		if err != nil {
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

	data, err := r.request(config.Cfg.MsgnetURL, request)
	if err != nil {
		return nil, err
	}

	serverResponse := new(Resp)
	err = json.Unmarshal(data, &serverResponse)
	if err != nil {
		return nil, err
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
	return nodeCert, nodeCert.Error.(error)
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

	return version.Result.Version, version.Error.(error)
}
