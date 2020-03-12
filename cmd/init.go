package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/corpulent/nuxx/pkg"
	"github.com/corpulent/nuxx/pkg/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a nuxx.yaml authentication file.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cliToken := viper.Get("token")

		if cliToken == "" {
			fmt.Println("Token is not set.")
			email, _ := common.PromptString("Your email.  Token activation link will be sent.")
			payloadString := fmt.Sprintf(`{"email": "%s"}`, email)
			json.Unmarshal([]byte(payloadString), &emailPayloadMap)
			jsonBytes, _ := json.MarshalIndent(emailPayloadMap, "", "  ")
			payload := string(jsonBytes)
			resp := pkg.PostRequest(getTokenEndpoint, payload)
			_ = json.Unmarshal([]byte(resp), &tokenRespData)
			fmt.Println(tokenRespData.RESP.API_ACCESS_KEY)
			apiAccessKey := tokenRespData.RESP.API_ACCESS_KEY

			if apiAccessKey != "" {
				viper.Set("token", apiAccessKey)
				if err := viper.WriteConfig(); err != nil {
					fmt.Println(err.Error())
				}

				fmt.Println("Your new cli access token was generated and set! Check your email for an activation link.")
			} else {
				fmt.Println("Oops.. something went wrong with setting cli access token.")
			}
		} else {
			fmt.Println("Your good to go.")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
