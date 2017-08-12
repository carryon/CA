package cmd

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
	"github.com/manucorporat/stats"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start deploy server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ConfigRuntime()
		StartWorkers()
		StartGin()
	},
}

var (
	mutexStats sync.RWMutex
	savedStats map[string]uint64
	messages   = stats.New()
	users      = stats.New()
	ips        = stats.New()
	port       string
)

type serverRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	Id     uint             `json:"id"`
}
type serverResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
	Id     uint        `json:"id"`
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&port, "port", "", "serve listening port.")
}

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}
func StartWorkers() {
	go statsWorker()
}

func statsWorker() {
	c := time.Tick(1 * time.Second)
	var lastMallocs uint64
	var lastFrees uint64
	for range c {
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)

		mutexStats.Lock()
		savedStats = map[string]uint64{
			"timestamp":  uint64(time.Now().Unix()),
			"HeapInuse":  stats.HeapInuse,
			"StackInuse": stats.StackInuse,
			"Mallocs":    (stats.Mallocs - lastMallocs),
			"Frees":      (stats.Frees - lastFrees),
			"Inbound":    uint64(messages.Get("inbound")),
			"Outbound":   uint64(messages.Get("outbound")),
			"Connected":  connectedUsers(),
		}
		lastMallocs = stats.Mallocs
		lastFrees = stats.Frees
		messages.Reset()
		mutexStats.Unlock()
	}
}
func connectedUsers() uint64 {
	connected := users.Get("connected") - users.Get("disconnected")
	if connected < 0 {
		return 0
	}
	return uint64(connected)
}

func StartGin() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(rateLimit, gin.Recovery())

	router.POST("/", func(c *gin.Context) {
		req := new(serverRequest)
		err := json.NewDecoder(c.Request.Body).Decode(req)

		res := new(serverResponse)
		res.Id = req.Id

		if err != nil {
			res.Error = err
			c.JSON(200, res)
		}

		var params [1]interface{}
		if req.Params != nil {
			var args interface{}
			params = [1]interface{}{args}
			if err = json.Unmarshal(*req.Params, &params); err != nil {
				res.Error = err
				c.JSON(200, res)
			}
		}

		var r interface{}
		var e error
		switch req.Method {
		case "msgnet-config":
			nodeType = "msg-net"
			r, e = nodesConfig(params)
		case "nodes-config":
			nodeType = "lcnd"
			r, e = nodesConfig(params)
		case "config-timestamp":
			nodeType = "lcnd"
			r, e = configTimestamp(params)
		case "msgnet-timestamp":
			nodeType = "msg-net"
			r, e = configTimestamp(params)
		case "lcnd-version":
			nodeType = "lcnd"
			r, e = nodeVersion()
		case "msgnet-version":
			nodeType = "msg-net"
			r, e = nodeVersion()
		default:
			e = errors.New("Invalid method")
		}
		res.Result = r
		if e != nil {
			res.Error = fmt.Sprintf("error message: %q", e)
		}
		c.JSON(200, res)
	})

	router.Run(":" + port)
}

func rateLimit(c *gin.Context) {
	ip := c.ClientIP()
	value := int(ips.Add(ip, 1))
	if value%50 == 0 {
		fmt.Printf("ip: %s, count: %d\n", ip, value)
	}
	if value >= 200 {
		if value%200 == 0 {
			fmt.Println("ip blocked")
		}
		c.Abort()
		c.String(503, "you were automatically banned :)")
	}
}

func nodesConfig(params interface{}) (interface{}, error) {
	aid := getAid(params)
	relationSlice, err := getRelationByAid(aid)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, v := range relationSlice {
		config, _ := deployDB.Get(nodeKey(v))
		var body interface{}

		if err := yaml.Unmarshal(config, &body); err != nil {
			return nil, err
		}
		var b map[interface{}]interface{} = body.(map[interface{}]interface{})
		b[interface{}("node_id")] = interface{}(v)
		if t, err := deployDB.Get(updateTimeKey(v)); err != nil {
			return nil, err
		} else {
			b[interface{}("update_time")] = interface{}(int64(binary.LittleEndian.Uint64(t)))
		}

		body = convert(b)
		if r, err := json.Marshal(body); err != nil {
			return nil, err
		} else {
			result = append(result, string(r))
		}
	}
	return result, nil
}

func configTimestamp(params interface{}) (interface{}, error) {
	aid := getAid(params)
	relationSlice, err := getRelationByAid(aid)
	if err != nil {
		return nil, err
	}
	var result []interface{}
	for _, v := range relationSlice {
		var r = make(map[string]interface{})
		updateTime, err := deployDB.Get(updateTimeKey(v))
		if err != nil {
			return nil, err
		}
		r["update_time"] = int64(binary.LittleEndian.Uint64(updateTime))
		r["node_id"] = v
		result = append(result, r)
	}
	return result, nil
}

func nodeVersion() (interface{}, error) {
	var result = make(map[string]string)
	r, err := deployDB.Get(versionKey())
	if err != nil {
		return nil, err
	}
	result["version"] = string(r)
	return result, nil
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

func getAid(params interface{}) string {
	tem := params.([1]interface{})
	return tem[0].(string)
}

func getRelationByAid(aid string) ([]string, error) {
	if len(aid) == 0 {
		return nil, errors.New("Invalid aid")
	}
	relation, err := deployDB.Get(relationKey(aid))
	if err != nil {
		return nil, err
	}

	return strings.Split(string(relation), ","), nil

}
