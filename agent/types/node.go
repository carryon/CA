package types

type NodeConfig struct {
	Net struct {
		MaxPeers       int           `json:"maxPeers"`
		BootstrapNodes []interface{} `json:"bootstrapNodes"`
		ListenAddr     string        `json:"listenAddr"`
		Privatekey     string        `json:"privatekey"`

		Msgnet struct {
			RouteAddress []string `json:"routeAddress"`
		} `json:"msgnet"`
	} `json:"net"`

	Log struct {
		Level string `json:"level"`
	} `json:"log"`

	Jrpc struct {
		Enabled bool   `json:"enabled"`
		Port    string `json:"port"`
	} `json:"jrpc"`

	Blockchain struct {
		ID         string `json:"id"`
		Datadir    string `json:"datadir"`
		Cpuprofile string `json:"cpuprofile"`
		ProfPort   string `json:"profPort"`
	} `json:"blockchain"`

	Issueaddr struct {
		Addr []string `json:"addr"`
	} `json:"issueaddr"`

	Consensus struct {
		Plugin string `json:"plugin"`

		Noops struct {
			BlockSize     int    `json:"blockSize"`
			BlockInterval string `json:"blockInterval"`
		} `json:"noops"`

		Lbft struct {
			ID                   string `json:"id"`
			N                    int    `json:"N"`
			Q                    int    `json:"Q"`
			K                    int    `json:"K"`
			BlockSize            int    `json:"blockSize"`
			BlockTimeout         string `json:"blockTimeout"`
			BlockInterval        string `json:"blockInterval"`
			BlockDelay           string `json:"blockDelay"`
			ViewChange           string `json:"viewChange"`
			ResendViewChange     string `json:"resendViewChange"`
			ViewChangePeriod     string `json:"viewChangePeriod"`
			NullRequest          string `json:"nullRequest"`
			BufferSize           int    `json:"bufferSize"`
			MaxConcurrentNumFrom int    `json:"maxConcurrentNumFrom"`
			MaxConcurrentNumTo   int    `json:"maxConcurrentNumTo"`
		} `json:"lbft"`
	} `json:"consensus"`

	Vm struct {
		MaxMem                     int    `json:"maxMem"`
		RegistrySize               int    `json:"registrySize"`
		CallStackSize              int    `json:"callStackSize"`
		ExecLimitStackDepth        int    `json:"execLimitStackDepth"`
		ExecLimitMaxOpcodeCount    int    `json:"execLimitMaxOpcodeCount"`
		ExecLimitMaxRunTime        int    `json:"execLimitMaxRunTime"`
		ExecLimitMaxScriptSize     int    `json:"execLimitMaxScriptSize"`
		ExecLimitMaxStateValueSize int    `json:"execLimitMaxStateValueSize"`
		ExecLimitMaxStateItemCount int    `json:"execLimitMaxStateItemCount"`
		ExecLimitMaxStateKeyLength int    `josn:"execLimitMaxStateKeyLength"`
		LuaVMExeFilePath           string `json:"luaVMExeFilePath"`
		JsVMExeFilePath            string `json:"jsVMExeFilePath"`
	} `json:"vm"`

	NodeID     string `json:"node_id"`
	UpdateTime int    `json:"update_time"`
}
