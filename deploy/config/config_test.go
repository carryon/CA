package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg, err := loadConfig("test.yaml")
	fmt.Println(cfg, *cfg.DBConfig, err)
}
