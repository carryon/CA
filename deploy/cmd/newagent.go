package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	aid    string
	remark string
	prefix string = "aid_"
)

// newagentCmd represents the newagent command
var newagentCmd = &cobra.Command{
	Use:   "newagent",
	Short: "Add a new agent",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(aid) == 0 {
			fmt.Println("Invalid aid")
			os.Exit(-1)
		}
		if len(remark) == 0 {
			fmt.Println("Invalid remark")
		}
		// aid_idval => remark
		err := deployDB.Put(aidKey(aid), []byte(remark))

		if err != nil {
			fmt.Printf("Failed to store aid, %v", err)
			os.Exit(-1)
		}
		fmt.Println("Succeed in storing aid")
		os.Exit(0)

	},
}

func aidKey(aid string) []byte {
	return []byte(prefix + aid)
}

func init() {
	RootCmd.AddCommand(newagentCmd)

	newagentCmd.Flags().StringVar(&aid, "aid", "", "agent ID")
	newagentCmd.Flags().StringVar(&remark, "remark", "", "remark")
}
