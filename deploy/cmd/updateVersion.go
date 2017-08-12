package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version string
)

// updateVersion represents the lcnd command
var updateVersion = &cobra.Command{
	Use:   "updateVersion",
	Short: "update version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(version) == 0 {
			fmt.Println("Invalid version")
			os.Exit(-1)
		}

		if err := deployDB.Put(versionKey(), []byte(version)); err != nil {
			fmt.Printf("Failed to update %s version, %v", nodeType, err)
			os.Exit(-1)
		}

		fmt.Printf("Succeed in updating %s version\r\n", nodeType)
	},
}

func versionKey() []byte {
	return []byte(nodeType + "_version")
}

func init() {
	RootCmd.AddCommand(updateVersion)

	updateVersion.Flags().StringVar(&version, "version", "", "node version")
	updateVersion.Flags().StringVar(&nodeType, "nodeType", "", "node type, lcnd or msg-net")
}
