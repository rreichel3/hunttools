package cmd

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	root "github.com/rreichel3/hunttools/cmd/root_flags"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

func init() {
	redisUploadCmd.Flags().StringVarP(&UploadJsonInfile, "in", "i", "", "Infile generated by ht redis dump -j. If a directory, will try to upload whole directory")
	redisUploadCmd.MarkFlagRequired("in")

	RedisRootCmd.AddCommand(redisUploadCmd)
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

func uploadData(client *redis.Client, serviceName string, dataToUpload []RedisData) int {
	var keyCount = 0
	for _, redisData := range dataToUpload {
		var key = fmt.Sprintf("%s:%s:%s", serviceName, strconv.Itoa(redisData.Database), redisData.Key)
		keyCount++
		// Need to base64 decode the value
		value := fmt.Sprintf("%v", redisData.Value)
		if root.VerboseOutput {
			fmt.Printf("restoring key: %v\n", key)
		}
		dumpValue, _ := b64.StdEncoding.DecodeString(value)
		_, err := client.Restore(ctx, key, 0, string(dumpValue)).Result()
		if err != nil {
			fmt.Println(err)
			continue
		}

	}
	return keyCount
}

func getFilename(fullPath string) string {
	fileName := filepath.Base(fullPath)
	var extension = filepath.Ext(fileName)
	return fileName[0 : len(fileName)-len(extension)]
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
		defer rdb.Close()
		fi, err := os.Stat(UploadJsonInfile)
		if err != nil {
			fmt.Println(err)
			return err
		}
		var pathsToUpload = []string{}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			// do directory stuff
			// Iterate over the directory and get every json file
			files, err := os.ReadDir(UploadJsonInfile)
			if err != nil {
				fmt.Println(err)
				return err
			}
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				pathsToUpload = append(pathsToUpload, filepath.Join(UploadJsonInfile, file.Name()))
			}

		case mode.IsRegular():
			// do file stuff
			fmt.Println("file")
			pathsToUpload = append(pathsToUpload, UploadJsonInfile)
		}
		for _, path := range pathsToUpload {
			if root.VerboseOutput {
				fmt.Printf("Loading file from %s\n", path)
			}
			// Load JSON Infile
			serviceName := getFilename(path)
			dataToUpload, err := loadRedisData(path)
			if err != nil {
				fmt.Printf("Unable to load file: %s\n", path)
				return err
			}
			// For each key, add to database
			keyCount := uploadData(rdb, serviceName, dataToUpload)
			if root.VerboseOutput {
				fmt.Printf("loaded %d keys\n", keyCount)
			}
		}
		return nil

	},
}
