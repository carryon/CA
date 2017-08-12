package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	nodeIDs string
)

// relateCmd represents the relate command
var relateCmd = &cobra.Command{
	Use:   "relate",
	Short: "relate agent and node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(aid) == 0 {
			fmt.Printf("Invalid agent ID")
			os.Exit(-1)
		}
		if len(nodeIDs) == 0 {
			fmt.Printf("Invalid node IDs")
			os.Exit(-1)
		}
		// check nodeIDs
		nodeIDSlice := strings.Split(nodeIDs, ",")
		for _, v := range nodeIDSlice {
			_, err := deployDB.Get(nodeKey(v))
			if err != nil {
				fmt.Printf("Unknown node %s \r\n", v)
				os.Exit(-1)
			}
		}

		// check if agent id exists
		agent, err := deployDB.Get(aidKey(aid))
		if err != nil || len(agent) == 0 {
			fmt.Println("Unknown agent")
			os.Exit(-1)
		}

		// old relation
		var relationSlice []string
		relation, _ := deployDB.Get(relationKey(aid))
		if relation != nil {
			relationSlice = strings.Split(string(relation), ",")
			relationSlice = append(relationSlice, nodeIDSlice...)
		} else {
			relationSlice = nodeIDSlice
		}
		removeDuplicatesAndEmpty(&relationSlice)
		newRelation := strings.Join(relationSlice, ",")

		if err := deployDB.Put(relationKey(aid), []byte(newRelation)); err != nil {
			fmt.Printf("Failed to restore relation, %v", err)
			os.Exit(-1)
		}
		fmt.Println("Succeed in relating")
	},
}

func relationKey(aid string) []byte {
	return []byte(nodeType + "_relation_" + aid)
}
func removeDuplicatesAndEmpty(slice *[]string) {
	found := make(map[string]bool)
	total := 0
	for i, val := range *slice {
		if _, ok := found[val]; !ok {
			found[val] = true
			(*slice)[total] = (*slice)[i]
			total++
		}
	}
	*slice = (*slice)[:total]
}

func init() {
	RootCmd.AddCommand(relateCmd)

	relateCmd.Flags().StringVar(&aid, "aid", "", "agent ID")
	relateCmd.Flags().StringVar(&nodeIDs, "nodeIDs", "", "node ID")
	relateCmd.Flags().StringVar(&nodeType, "nodeType", "", "node type, lcnd or msg-net")
}
