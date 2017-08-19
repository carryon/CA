package types

type NodeConfig struct {
	Blockchain struct {
		Datadir string `json:"datadir"`
		ID      string `json:"id"`
	} `json:"blockchain"`
	Consensus struct {
		Lbft struct {
			K                int    `json:"K"`
			N                int    `json:"N"`
			Q                int    `json:"Q"`
			BlockDelay       string `json:"blockDelay"`
			BlockInterval    string `json:"blockInterval"`
			BlockSize        int    `json:"blockSize"`
			BlockTimeout     string `json:"blockTimeout"`
			BufferSize       int    `json:"bufferSize"`
			ID               string `json:"id"`
			MaxConcurrentNum int    `json:"maxConcurrentNum"`
			NullRequest      string `json:"nullRequest"`
			ResendViewChange string `json:"resendViewChange"`
			ViewChange       string `json:"viewChange"`
			ViewChangePeriod string `json:"viewChangePeriod"`
		} `json:"lbft"`
		Noops struct {
			BlockInterval string `json:"blockInterval"`
			BlockSize     int    `json:"blockSize"`
		} `json:"noops"`
		Plugin string `json:"plugin"`
	} `json:"consensus"`
	Issueaddr struct {
		Addr []string `json:"addr"`
	} `json:"issueaddr"`
	Jrpc struct {
		Enabled bool   `json:"enabled"`
		Port    string `json:"port"`
	} `json:"jrpc"`
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
	Net struct {
		BootstrapNodes []interface{} `json:"bootstrapNodes"`
		ListenAddr     string        `json:"listenAddr"`
		MaxPeers       int           `json:"maxPeers"`
		Msgnet         struct {
			RouteAddress []string `json:"routeAddress"`
		} `json:"msgnet"`
		Privatekey string `json:"privatekey"`
	} `json:"net"`
	NodeID     string `json:"node_id"`
	UpdateTime int    `json:"update_time"`
}
