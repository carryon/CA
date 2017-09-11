package config

import (
	"fmt"
	"path/filepath"

	"github.com/bocheninc/CA/agent/utils"
	"github.com/spf13/viper"
)

var (
	Cfg                  *Config
	defaultExecDirName   = "bin"
	defaultLcndDirName   = "lcnd"
	defaultMsgnetDirName = "msgnet"
	defaultConfigFile    = "agent.yaml"
	defaultLogLevel      = "debug"
	defaultLogFilename   = "agent.log"
)

type Config struct {
	ID           string
	DeployServer string
	BaseDir      string
	ExecDir      string
	LcndDir      string
	MsgNetDir    string
	MsgnetURL    string
	LcndURL      string
	LogLevel     string
	LogFile      string
}

func LoadConfig(cfgfile string) (*Config, error) {
	var err error
	conf := new(Config)

	viper.SetConfigFile(cfgfile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf(" viper.ReadInConfig error %s", err)
	}

	conf.readConfig()

	utils.OpenDir(conf.BaseDir)

	if conf.ExecDir, err = utils.OpenDir(filepath.Join(conf.BaseDir, defaultExecDirName)); err != nil {
		return nil, err
	}
	if conf.LcndDir, err = utils.OpenDir(filepath.Join(conf.BaseDir, defaultLcndDirName)); err != nil {
		return nil, err
	}
	if conf.MsgNetDir, err = utils.OpenDir(filepath.Join(conf.BaseDir, defaultMsgnetDirName)); err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) readConfig() {
	if id := viper.GetString("id"); id != "" {
		c.ID = id
	}

	if deploy := viper.GetString("server.deploy"); deploy != "" {
		c.DeployServer = deploy
	}
	if basedir := viper.GetString("basedir"); basedir != "" {
		c.BaseDir = basedir
	}
	if lcnd := viper.GetString("URL.lcnd"); lcnd != "" {
		c.LcndURL = lcnd
	}
	if msgnet := viper.GetString("URL.msgnet"); msgnet != "" {
		c.MsgnetURL = msgnet
	}

	if logLevel := viper.GetString("log.level"); logLevel != "" {
		c.LogLevel = logLevel
		c.LogFile = defaultLogFilename
	}

}
