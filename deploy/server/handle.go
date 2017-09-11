package server

import (
	"encoding/json"

	"fmt"

	yaml "gopkg.in/yaml.v2"
)

func msgnetsConfig(params []string, list *List) (interface{}, error) {
	//todo msgnet config

	return nil, nil
}

func nodesConfig(params []string, list *List) (interface{}, error) {
	aid := params[0]
	nodes, ok := list.NodeList[aid]

	fmt.Println("agentID: ", aid)
	if !ok {
		return nil, fmt.Errorf("not found by agent id : %s", aid)
	}

	var result []string
	for _, v := range nodes {
		var body interface{}
		if err := yaml.Unmarshal([]byte(v.ConfigFile), &body); err != nil {
			return nil, err
		}

		b := body.(map[interface{}]interface{})

		b[interface{}("blockchain")].(map[interface{}]interface{})[interface{}("chainId")] = interface{}(v.ChainID)
		b[interface{}("blockchain")].(map[interface{}]interface{})[interface{}("nodeId")] = interface{}(v.NodeID)

		body = convert(b)
		r, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		result = append(result, string(r))
	}
	//log.Debugf("nodes Config result : %v len: %d", result, len(nodes))
	return result, nil
}

func lcndConfigTimestamp(params []string, list *List) (interface{}, error) {
	aid := params[0]
	nodes := list.NodeList[aid]
	var result []interface{}
	for _, v := range nodes {
		var r = make(map[string]interface{})
		//todo defult update_time
		r["update_time"] = int64(123456)
		//r["update_time"] = v.Updated.Unix()
		r["node_id"] = v.NodeID
		result = append(result, r)
	}
	return result, nil
}

func msgnetConfigTimestamp(params []string, list *List) (interface{}, error) {
	//todo
	return nil, nil
}

func nodeVersion(list *List) (interface{}, error) {
	var result = make(map[string]string)
	result["version"] = "v0.8.8"
	return result, nil
}

func nodeCert(params []string, ca *Ca) (interface{}, error) {
	chainID := params[0]
	nodeID := params[1]
	key := params[2]

	certificate, err := ca.GetCert(chainID, nodeID, key)
	if err != nil {
		return nil, err
	}

	rootCertificate, err := ca.GetRootCertificate()
	if err != nil {
		return nil, err
	}

	return []string{string(certificate), string(rootCertificate)}, nil
}

func msgnetVersion(list *List) (interface{}, error) {
	//todo
	return nil, nil
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
