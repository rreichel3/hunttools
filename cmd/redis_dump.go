package cmd

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

func init() {
	redisDumpCmd.Flags().IntVarP(&MaxDump, "count", "n", -1, "Max number to dump. Defaults to all. May dump a few more depending on how much comes from redis.")
	redisDumpCmd.Flags().BoolVarP(&DumpJson, "json", "j", false, "Dump as json")
	redisDumpCmd.Flags().BoolVarP(&DumpValues, "values", "", false, "Dump values with keys")
	redisRootCmd.AddCommand(redisDumpCmd)
}

type RedisData struct {
	Hostname string `json:"hostname"`
	Database int    `json:"database"`
	Key      string `json:"key"`
	Value    string `json:"value"`
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

		keyspace, _ := rdb.Info(ctx, "keyspace").Result()
		re := regexp.MustCompile(`db(\d+)`)
		matches := re.FindAllStringSubmatch(keyspace, -1)
		dbs := make(map[int]bool, 1)
		for _, match := range matches {
			dbNum, _ := strconv.Atoi(match[1])
			dbs[dbNum] = true
		}
		if VerboseOutput {
			fmt.Printf("Dumping %d keyspace(s)..\n", len(dbs))
			fmt.Printf("%v", dbs)
		}

		// A map of databases to keys
		var allKeys = make(map[int][]string, 0)

		for dbNum, _ := range dbs {
			rdb := redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", RedisHostname, RedisPort),
				Password: RedisPassword, // no password set
				DB:       dbNum,         // use default DB
			})
			var keysForDb []string = make([]string, 0)

			for {
				var keys []string
				var err error
				keys, cursor, err = rdb.Scan(ctx, cursor, "*", 1000).Result()
				if err != nil {
					panic(err)
				}
				keysForDb = append(keysForDb, keys...)

				n += len(keys)
				if cursor == 0 || (MaxDump != -1 && n >= MaxDump) {
					break
				}
			}
			allKeys[dbNum] = keysForDb
		}

		if VerboseOutput {
			fmt.Printf("found %d keys\n", n)
			fmt.Printf("found %v keys\n", allKeys)
		}

		if DumpValues {
			var fullList = make([]RedisData, len(allKeys))
			for dbNum, keys := range allKeys {
				rdb = redis.NewClient(&redis.Options{
					Addr:     fmt.Sprintf("%s:%d", RedisHostname, RedisPort),
					Password: RedisPassword, // no password set
					DB:       dbNum,         // use default DB
				})
				// limiter := redis_rate.NewLimiter(rdb)
				// _, err := limiter.Allow(ctx, "*", redis_rate.PerSecond(1))
				// if err != nil {
				// panic(err)
				// }

				for _, key := range keys {
					time.Sleep(time.Millisecond * 100)
					val, err := rdb.Dump(ctx, key).Result()
					if err != nil {
						if VerboseOutput {
							fmt.Printf("Error encountered getting key: %s\n", key)
							fmt.Printf("Error: %v\n", err)
						}
						continue
					}
					if VerboseOutput {
						fmt.Printf("Got key: %s\n", key)
					}
					sEnc := b64.StdEncoding.EncodeToString([]byte(val))
					fullList = append(fullList, RedisData{Hostname: RedisHostname, Database: dbNum, Key: key, Value: sEnc})
				}
			}
			idx := 1
			if DumpJson {
				jsonStr, err := json.Marshal(fullList)
				if err != nil {
					fmt.Printf("Error encountered marshaling json: %s\n", err)
					return err
				}
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
			for dbNum, key := range allKeys {
				if VerboseOutput {
					fmt.Printf("%d: (DB%d) %s\n", idx, dbNum, key)
				} else {
					fmt.Printf("(DB%d) %s\n", dbNum, key)
				}
				idx++
			}
		}

		return nil

	},
}
