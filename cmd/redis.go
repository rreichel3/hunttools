package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	redisRootCmd.Flags().StringVarP(&RedisHostname, "host", "H", "localhost", "Redis hostname. Defaults to locahost")
	redisRootCmd.Flags().IntVarP(&RedisPort, "port", "P", 6379, "Redis Port. Defaults to 6379")
	redisRootCmd.Flags().StringVarP(&RedisPassword, "password", "p", "", "Redis password. Defaults to empty.")
	redisRootCmd.Flags().IntVarP(&RedisDB, "database", "d", 0, "Redis database. Defaults to 0.")
	rootCmd.AddCommand(redisRootCmd)
}

var ctx = context.Background()

var RedisHostname string
var RedisPort int
var RedisPassword string
var RedisDB int
var redisRootCmd = &cobra.Command{
	Use:   "redis",
	Short: "Commands to get redis ",
	Long:  `Dump all keys and values in Redis.. Safely. Uses SCAN.`,
}
