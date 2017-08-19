package http

type VersionAPI struct {
	ID     int `json:"id"`
	Result struct {
		Version string `json:"version"`
	} `json:"result"`
	Error interface{} `json:"error"`
}

type DownloadAPI struct {
	ID     int         `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

type TimestampAPI struct {
	Result []struct {
		NodeID     string `json:"node_id"`
		UpdateTime int    `json:"update_time"`
	} `json:"result"`
	Error interface{} `json:"error"`
	ID    int         `json:"id"`
}
