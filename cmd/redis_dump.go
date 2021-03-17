package cmd

import (
	"fmt"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

func init() {
	redisDumpCmd.Flags().IntVarP(&MaxDump, "count", "n", -1, "Max number to dump. Defaults to all. May dump a few more depending on how much comes from redis.")
	redisDumpCmd.Flags().BoolVarP(&DumpJson, "json", "j", false, "Dump as json")
	redisRootCmd.AddCommand(redisDumpCmd)
}

var MaxDump int 
var DumpJson bool
var redisDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump all keys and values in Redis.. Safely",
	Long:  `Dump all keys and values in Redis.. Safely. Uses SCAN.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", RedisHostname, RedisPort),
			Password: RedisPassword, // no password set
			DB:       0,  // use default DB
		})
		var cursor uint64
		var n int
		var allKeys = make(map[string]string)
		for {
			var keys []string
			var err error
			keys, cursor, err = rdb.Scan(ctx, cursor, "*", 10).Result()
			if err != nil {
				panic(err)
			}
			for _, key := range keys {
				allKeys[key] = ""
			}
			n += len(keys)
			if cursor == 0 || (MaxDump != -1 && n >= MaxDump) {
				break
			}
		}
		if VerboseOutput {
			fmt.Printf("found %d keys\n", n)
		}
		idx := 1
		for key, _ := range allKeys {
			val, err := rdb.Get(ctx, key).Result()
			if err != nil {
				fmt.Printf("Error encountered getting key: %s\n", key)
				continue
			}
			allKeys[key] = val
			idx++
		}
		idx = 1
		if DumpJson {
			jsonStr, _ := json.Marshal(allKeys)
			fmt.Println(string(jsonStr))
		} else {
			for k, v := range allKeys {
				if VerboseOutput {
					fmt.Printf("%d: %s => %s\n", idx, k, v)
				} else {
					fmt.Printf("%s => %s\n", k, v)
				}
				idx++
			}
		}
		
		return nil

	},
}

