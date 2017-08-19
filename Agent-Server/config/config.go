package config

type Config struct {
	ID              string `yaml:"id"`
	CaServer        string `yaml:"ca_server"`
	DeployServer    string `yaml:"deploy_server"`
	BaseDir         string `yaml:"base_dir"`
	AgentCertDir    string `yaml:"agent_cert_dir"`
	L0ExecFile      string `yaml:"l0_exec_file"`
	MsgnetConfigDir string `yaml:"msgnet_config_dir"`
	MsgnetExecFile  string `yaml:"msgnet_exec_file"`
	MsgnetUrl       string `yaml:"msgnet_url"`
	LcndUrl         string `yaml:"lcnd_url"`
}
