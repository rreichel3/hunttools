package cmd

import (
	"fmt"

	"github.com/rreichel3/hunttools/cmd/utils"

	"github.com/spf13/cobra"
)

func init() {
	pingCmd.Flags().StringVarP(&PingDestinationsPath, "in", "i", "", "Newline delimited file of ping destinations")
	pingCmd.Flags().BoolVarP(&PingAliveFlag, "alive", "a", false, "Only output if the host is alive")
	rootCmd.AddCommand(pingCmd)
}

var PingDestinationsPath string
var PingAliveFlag bool

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping a list of hosts",
	Long:  `Runs the ping command against a list of hosts`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var pingDestinations, err = utils.ReadFileToList(PingDestinationsPath)
		if err != nil {
			return err
		}

		for _, each_ln := range pingDestinations {
			fmt.Println(each_ln)
		}
		return nil

	},
}
