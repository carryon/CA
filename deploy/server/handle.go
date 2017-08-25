package server

import (
	"encoding/json"

	"fmt"

	yaml "gopkg.in/yaml.v2"
)

func msgnetsConfig(params interface{}, list *List) (interface{}, error) {
	//todo msgnet config
	// aid := getID(params)
	// _ = list.AgentList[aid]

	return nil, nil
}

func nodesConfig(params interface{}, list *List) (interface{}, error) {
	aid := getID(params)
	nodes, ok := list.NodeList[aid]
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

		b[interface{}("node_id")] = interface{}(v.NodeID)
		b[interface{}("update_time")] = interface{}(v.Updated.Unix())

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

func lcndConfigTimestamp(params interface{}, list *List) (interface{}, error) {
	aid := getID(params)
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

func msgnetConfigTimestamp(params interface{}, list *List) (interface{}, error) {
	//todo
	return nil, nil
}

func nodeVersion(list *List) (interface{}, error) {
	var result = make(map[string]string)
	result["version"] = "0.88"
	return result, nil
}

func msgnetVersion(list *List) (interface{}, error) {
	//todo
	return nil, nil
}

func getID(params interface{}) string {
	tem := params.([1]interface{})
	return tem[0].(string)
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
