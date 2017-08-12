package cmd

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	config   string
	nodeID   string
	nodeType string
)

// newconfigCmd represents the newconfig command
var newconfigCmd = &cobra.Command{
	Use:   "newconfig",
	Short: "create the specified node's or msg-net config from the config file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(config) == 0 {
			fmt.Printf("Invalid config")
			os.Exit(-1)
		}
		if len(nodeID) == 0 {
			fmt.Printf("Invalid node ID")
			os.Exit(-1)
		}
		if nodeType != "msg-net" {
			nodeType = "lcnd"
		}
		pwd, _ := os.Getwd()
		if _, err := os.Stat(pwd + "/" + config); err != nil {
			fmt.Printf("Invalid config, %v", err)
			os.Exit(-1)
		}

		// read config
		data, err := ioutil.ReadFile(pwd + "/" + config)
		if err != nil {
			fmt.Printf("Failed to read config, %v", err)
			os.Exit(-1)
		}
		/*
			var dm map[string]interface{}
			yaml.Unmarshal(data, &dm)
		*/
		// check config

		// store config
		if err = deployDB.Put(nodeKey(nodeID), data); err != nil {
			fmt.Printf("Failed to store node config, %v", err)
			os.Exit(-1)
		}

		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(time.Now().Unix()))
		if err = deployDB.Put(updateTimeKey(nodeID), b); err != nil {
			fmt.Printf("Failed to update modification time, %v", err)
			os.Exit(-1)
		}
		fmt.Println("Succeed in create a node config")
		os.Exit(0)
	},
}

func nodeKey(ID string) []byte {
	return []byte(nodeType + "_config_" + ID)
}
func updateTimeKey(ID string) []byte {
	return []byte(nodeType + "_update_" + ID)
}

func init() {
	RootCmd.AddCommand(newconfigCmd)

	newconfigCmd.Flags().StringVar(&config, "config", "", "Path to the node config")
	newconfigCmd.Flags().StringVar(&nodeID, "nodeID", "", "Node ID")
	newconfigCmd.Flags().StringVar(&nodeType, "nodeType", "lcnd", "lcnd or msg-net")
}
