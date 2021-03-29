package heroku

import (
	"encoding/json"
	"fmt"

	heroku "github.com/heroku/heroku-go/v5"
	"github.com/spf13/cobra"
	"os"
)

func init() {

	herokuRootCommand.AddCommand(herokuDumpAllCmd)
}

type HerokuSecret struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type HerokuApp struct {
	Name    string         `json:"name"`
	ID      string         `json:"id"`
	Secrets []HerokuSecret `json:"secrets"`
}

var herokuDumpAllCmd = &cobra.Command{
	Use:   "dump-all",
	Short: "Dump all secrets for all apps",
	Long:  `Dump all secrets for all apps`,
	RunE: func(cmd *cobra.Command, args []string) error {
		bearer, ok := os.LookupEnv("HEROKU_API_KEY")
		if !ok {
			fmt.Println("You need to set the HEROKU_API_KEY environment variable.\n")
			return nil
		}

		heroku.DefaultTransport.BearerToken = bearer

		h := heroku.NewService(heroku.DefaultClient)
		apps, err := h.AppList(ctx, nil)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		appsDump := make([]*HerokuApp, 0)
		for _, app := range apps {
			println(app.Name)
			secrets, err := h.ConfigVarInfoForApp(ctx, app.ID)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			secretsList := []HerokuSecret{}
			for key, secret := range secrets {
				secretsList = append(secretsList, HerokuSecret{key, *secret})
			}
			herokuApp := HerokuApp{app.Name, app.ID, secretsList}
			appsDump = append(appsDump, &herokuApp)
		}
		jsonStr, err := json.Marshal(appsDump)
		if err != nil {
			fmt.Printf("Error encountered marshaling json: %s\n", err)
			return err
		}
		fmt.Println(string(jsonStr))

		return nil

	},
}
