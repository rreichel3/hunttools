package cmd

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

func init() {
	redisUploadCmd.Flags().StringVarP(&UploadJsonInfile, "infile", "i", "", "Infile generated by ht redis dump -j")
	redisUploadCmd.MarkFlagRequired("infile")

	redisRootCmd.AddCommand(redisUploadCmd)
}

func loadRedisData(sourceFilePath string) ([]RedisData, error) {
	jsonFile, err := os.Open(sourceFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []RedisData
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}

var UploadJsonInfile string
var redisUploadCmd = &cobra.Command{
	Use:   "restore",
	Short: "Upload the results of ht redis dump",
	Long:  `Upload the results of ht redis dump. Uses RESTORE`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", RedisHostname, RedisPort),
			Password: RedisPassword, // no password set
			DB:       RedisDB,       // use default DB
		})

		// Load JSON Infile
		dataToUpload, err := loadRedisData(UploadJsonInfile)
		if err != nil {
			fmt.Printf("Unable to load file: %s\n", UploadJsonInfile)
			return err
		}
		// For each key, add to database
		for _, redisData := range dataToUpload {
			var key = fmt.Sprintf("%s:%s", redisData.Database, redisData.Key)
			// Need to base64 decode the value
			value := fmt.Sprintf("%v", redisData.Value)
			dumpValue, _ := b64.StdEncoding.DecodeString(value)
			_, err := rdb.Restore(ctx, key, 0, string(dumpValue)).Result()
			if err != nil {
				fmt.Println(err)
				continue
			}

		}

		return nil

	},
}
