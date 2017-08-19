package types

type CertConfig struct {
	Cert struct {
		KeyPath string `yaml:"keyPath"`
		CrtPath string `yaml:"crtPath"`
		CaPath  string `yaml:"caPath"`
	}
}

type CAConfig struct {
	CA struct {
		Enabled bool `yaml:"enabled"`
	}
}

type NodeCert struct {
	ID      int `json:"id"`
	Results struct {
		NodeID   string `json:"node_id"`
		AgentKey string `json:"agent_key"`
		AgentCrt string `json:"agent_crt"`
		RootCrt  string `json:"root_crt"`
	} `json:"results"`
	Error interface{} `json:"error"`
}
