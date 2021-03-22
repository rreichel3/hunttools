package cmd

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

func init() {
	redisDumpCmd.Flags().IntVarP(&MaxDump, "count", "n", -1, "Max number to dump. Defaults to all. May dump a few more depending on how much comes from redis.")
	redisDumpCmd.Flags().BoolVarP(&DumpJson, "json", "j", false, "Dump as json")
	redisDumpCmd.Flags().BoolVarP(&DumpValues, "values", "", false, "Dump values with keys")
	redisRootCmd.AddCommand(redisDumpCmd)
}

var MaxDump int
var DumpJson bool
var DumpValues bool
var redisDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump all keys and values in Redis.. Safely",
	Long:  `Dump all keys and values in Redis.. Safely. Uses DUMP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", RedisHostname, RedisPort),
			Password: RedisPassword, // no password set
			DB:       RedisDB,       // use default DB
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
		if DumpValues {
			for key, _ := range allKeys {
				val, err := rdb.Dump(ctx, key).Result()
				if err != nil {
					fmt.Printf("Error encountered getting key: %s\n", key)
					continue
				}
				sEnc := b64.StdEncoding.EncodeToString([]byte(val))
				allKeys[key] = sEnc
			}
			idx := 1
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
		} else {
			idx := 1
			for k, _ := range allKeys {
				if VerboseOutput {
					fmt.Printf("%d: %s\n", idx, k)
				} else {
					fmt.Printf("%s\n", k)
				}
				idx++
			}
		}

		return nil

	},
}
