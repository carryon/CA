package cmd

import (
	"fmt"
	"os"

	"github.com/bocheninc/CA/agent/config"
	s "github.com/bocheninc/CA/agent/server"
	"github.com/bocheninc/L0/components/log"
	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "agent",
	Short: "A brief description of your application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		config.Cfg = cfg

		log.New(config.Cfg.LogFile)
		log.SetLevel(config.Cfg.LogLevel)

		log.Info("load config : ", *config.Cfg)
		agentServer := s.NewServer()
		agentServer.Start()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is agent.yaml)")
}
