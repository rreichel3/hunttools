package cmd

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	root "github.com/rreichel3/hunttools/cmd/root_flags"
	"github.com/rreichel3/hunttools/cmd/utils"

	"github.com/spf13/cobra"
)

func init() {
	redisPingCmd.Flags().StringVarP(&RedisPingInfile, "infile", "i", "", "Newline delimited file of IPs for which to fetch their hostname")
	redisPingCmd.MarkFlagRequired("infile")

	RedisRootCmd.AddCommand(redisPingCmd)
}

var RedisPingInfile string
var redisPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping server",
	Long:  `Try redis ping`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if root.VerboseOutput {
			fmt.Println("Loading in: " + RedisPingInfile)
		}
		var redisHosts, err = utils.ReadFileToList(RedisPingInfile)
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, redisHost := range redisHosts {
			rdb := redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", redisHost, RedisPort),
				Password: RedisPassword, // no password set
				DB:       RedisDB,       // use default DB
			})
			_, err := rdb.Ping(ctx).Result()
			if err == nil {
				fmt.Println(redisHost)
			} else {
				if root.VerboseOutput {
					fmt.Printf("Error encountered for %s: %v\n", redisHost, err)
				}
			}
		}

		return nil

	},
}
