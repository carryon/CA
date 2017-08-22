package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bocheninc/CA/deploy/components/log"
	"github.com/gin-gonic/gin"
)

type serverRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	ID     uint             `json:"id"`
}

type serverResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
	ID     uint        `json:"id"`
}

type Router struct {
	port   string
	ticker *time.Ticker
	engine *gin.Engine
	list   *List
	stats  *Stats
}

func NewRouter(list *List, port string) *Router {
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	return &Router{
		port:   port,
		ticker: time.NewTicker(5 * time.Second),
		engine: gin.Default(),
		list:   list,
		stats:  NewStats(),
	}
}

func (r *Router) start() {
	log.Info("router start ...")
	go r.eventLoop()
	go r.stats.statsWorker()

	r.engine.Use(r.rateLimit, gin.Recovery())

	r.engine.POST("/", r.handle)
	r.engine.Run(":" + r.port)
}

func (r *Router) eventLoop() {
	for {
		select {
		case <-r.ticker.C:
			r.list.UpdateNodeList()
			log.Debugf("update node list: %v ,update agent list: %v", r.list.NodeList, r.list.AgentList)
		}
	}
}

func (r *Router) rateLimit(c *gin.Context) {
	ip := c.ClientIP()
	value := int(r.stats.Ips.Add(ip, 1))
	if value%50 == 0 {
		log.Debugf("ip: %s, count: %d\n", ip, value)
	}
	if value >= 200 {
		if value%200 == 0 {
			log.Warnf("ip: %s blocked", ip)
		}
		c.Abort()
		c.String(503, "you were automatically banned :)")
	}
}

func (r *Router) handle(c *gin.Context) {
	req := new(serverRequest)
	err := json.NewDecoder(c.Request.Body).Decode(req)

	res := new(serverResponse)
	res.ID = req.ID

	if err != nil {
		res.Error = err.Error()
		c.JSON(200, res)
	}

	var params [1]interface{}

	if req.Params != nil {
		var args interface{}
		params = [1]interface{}{args}
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			res.Error = err.Error()
			c.JSON(200, res)
		}
	}

	var result interface{}
	log.Debug("agent request: ", *req)
	switch req.Method {
	case "msgnet-config":
		result, err = msgnetsConfig(params, r.list)
	case "nodes-config":
		result, err = nodesConfig(params, r.list)
	// case "lcnd-timestamp":
	// 	result, err = lcndConfigTimestamp(params, r.list)
	// case "msgnet-timestamp":
	// 	result, err = msgnetConfigTimestamp(params, r.list)
	case "lcnd-version":
		result, err = nodeVersion(r.list)
	case "msgnet-version":
		result, err = msgnetVersion(r.list)
	default:
		err = errors.New("Invalid method")
	}
	res.Result = result
	if err != nil {
		res.Error = fmt.Errorf("error message: %s", err).Error()
	}

	c.JSON(200, res)
}
