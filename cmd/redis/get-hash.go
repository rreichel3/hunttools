package cmd

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/go-redis/redis/v8"
	root "github.com/rreichel3/hunttools/cmd/root_flags"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func init() {
	redisGetHashCmd.Flags().StringVarP(&RedisKey, "key", "r", "", "Redis key to hash")
	redisGetHashCmd.MarkFlagRequired("redis-key")

	redisGetHashCmd.Flags().StringVarP(&HashKey, "field", "f", "", "Field in hash to fetch")
	redisGetHashCmd.MarkFlagRequired("hash-key")

	redisGetHashCmd.Flags().BoolVarP(&DecompressZlib, "zlib-decompress", "z", false, "Decompress value as compressed zlib")
	RedisRootCmd.AddCommand(redisGetHashCmd)
}

var RedisKey string
var HashKey string
var DecompressZlib bool
var redisGetHashCmd = &cobra.Command{
	Use:   "get-hash",
	Short: "Get specific hash value",
	Long:  `Get specific hash value. Can also decompress the value.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", RedisHostname, RedisPort),
			Password: RedisPassword, // no password set
			DB:       0,             // use default DB
		})
		if root.VerboseOutput {
			fmt.Println("Fetching key and value")
		}
		value, err := rdb.HGet(ctx, RedisKey, HashKey).Result()
		if err != nil {
			fmt.Println(err)
			return err
		}
		if DecompressZlib {
			if root.VerboseOutput {
				fmt.Println("Decompressing")
			}
			var b = bytes.NewBufferString(value)
			r, err := zlib.NewReader(b)
			if err != nil {
				fmt.Println(err)
				return err
			}
			io.Copy(os.Stdout, r)
			r.Close()
		} else {
			fmt.Println(value)
		}

		return nil

	},
}
