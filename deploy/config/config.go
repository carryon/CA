package config

import (
	"path/filepath"

	"github.com/bocheninc/CA/deploy/components/db"
	"github.com/bocheninc/CA/deploy/components/log"
	"github.com/bocheninc/CA/deploy/components/utils"
	"github.com/spf13/viper"
)

const (
	defaultLogLevel       = "debug"
	defaultLogFilename    = "deplog.log"
	defaultLogDirname     = "log"
	defaultConfigFilename = "deploy.yaml"
	defaultPort           = "8888"
)

var (
	defaultConfig = &Config{
		DBConfig: db.DefaultMysqldbConfig(),
		LogLevel: defaultLogLevel,
		LogFile:  defaultLogFilename,
		LogDir:   defaultLogDirname,
		Port:     defaultPort,
	}
)

type Config struct {
	LogLevel        string
	LogFile         string
	LogDir          string
	Port            string
	RouterAddresses []string
	DBConfig        *db.MysqldbConfig
}

// New returns a config according the config file
func NewConfig(cfgFile string) (cfg *Config, err error) {
	return loadConfig(cfgFile)
}

func loadConfig(cfgFile string) (conf *Config, err error) {
	var cfg *Config

	cfg = defaultConfig

	if cfgFile != "" {
		if utils.FileExist(cfgFile) {
			viper.SetConfigFile(cfgFile)
		}
		if err := viper.ReadInConfig(); err != nil {
			log.Warnf("no config file, run as default config! viper.ReadInConfig error %s", err)
		}
		cfg.LogDir, err = utils.OpenDir(defaultLogDirname)
		if err != nil {
			log.Error(err)
		}

		cfg.readLogConfig()
		cfg.readDBConfig()
		cfg.readPort()
		cfg.readMsgnetAddresses()
		return cfg, nil
	}

	return cfg, nil
}

func (cfg *Config) readLogConfig() {
	if logLevel := viper.GetString("log.level"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
	if logFile := filepath.Join(cfg.LogDir, defaultLogFilename); logFile != "" {
		cfg.LogFile = logFile
	}
}

func (cfg *Config) readDBConfig() {
	if dbName := viper.GetString("db.name"); dbName != "" {
		cfg.DBConfig.Name = dbName
	}
	if dbUser := viper.GetString("db.user"); dbUser != "" {
		cfg.DBConfig.User = dbUser
	}
	if dbPWD := viper.GetString("db.password"); dbPWD != "" {
		cfg.DBConfig.PWD = dbPWD
	}
	if dbHost := viper.GetString("db.host"); dbHost != "" {
		cfg.DBConfig.Host = dbHost
	}
	if dbPort := viper.GetString("db.port"); dbPort != "" {
		cfg.DBConfig.Port = dbPort
	}
	if dbZone := viper.GetString("db.zpne"); dbZone != "" {
		cfg.DBConfig.Zone = dbZone
	}
}

func (cfg *Config) readPort() {
	if port := viper.GetString("listen.port"); port != "" {
		cfg.Port = port
	}
}

func (cfg *Config) readMsgnetAddresses() {
	if routerAddresses := viper.GetStringSlice("msgnet.addresses"); routerAddresses != nil {
		cfg.RouterAddresses = routerAddresses
	}
}
