package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	RedisRootCmd.PersistentFlags().StringVarP(&RedisHostname, "host", "H", "localhost", "Redis hostname. Defaults to locahost")
	RedisRootCmd.PersistentFlags().IntVarP(&RedisPort, "port", "P", 6379, "Redis Port. Defaults to 6379")
	RedisRootCmd.PersistentFlags().StringVarP(&RedisPassword, "password", "p", "", "Redis password. Defaults to empty.")
	RedisRootCmd.PersistentFlags().IntVarP(&RedisDB, "database", "d", 0, "Redis db. Defaults to 0.")
}

var ctx = context.Background()

var RedisHostname string
var RedisPort int
var RedisPassword string
var RedisDB int
var RedisRootCmd = &cobra.Command{
	Use:   "redis",
	Short: "Commands to get redis ",
	Long:  `Dump all keys and values in Redis.. Safely. Uses SCAN.`,
}
