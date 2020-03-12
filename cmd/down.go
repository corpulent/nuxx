package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/corpulent/nuxx/pkg"
	"github.com/corpulent/nuxx/pkg/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/thedevsaddam/gojsonq"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Bring down a deployment.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(projectConfFile); os.IsNotExist(err) {
			d := color.New(color.FgCyan, color.Bold).SprintFunc()
			fmt.Printf(`Project configuration file does not exist! Create it by running: %s`, d("nuxx compose"))
			os.Exit(1)
		}

		cliToken := common.GetToken()
		file, _ := ioutil.ReadFile(projectConfFile)
		json.Unmarshal([]byte(file), &configStructureMap)
		existingData = gojsonq.New().File(projectConfFile)
		projectName := existingData.Find("project_name")

		downPathEndpoint := fmt.Sprintf(`%s%s`, downPathEndpoint, cliToken)
		downJsonStr := fmt.Sprintf(`{"project_name": "%s"}`, projectName)
		resp := pkg.DownRequest(downPathEndpoint, downJsonStr)

		printAllResponses(resp)
	},
}

func printAllResponses(resp *pkg.DownResponse) {
	respResp := resp.Resp
	p := fmt.Println

	if len(respResp.Notices) > 0 {
		p(``)
		for _, v := range respResp.Notices {
			p(v)
		}
		p(``)
	}

	if len(respResp.Warnings) > 0 {
		p(``)
		for _, v := range respResp.Warnings {
			p(v)
		}
		p(``)
	}

	if len(respResp.Errors) > 0 {
		p(``)
		for _, v := range respResp.Errors {
			p(v)
		}
		p(``)
	}
}

func init() {
	rootCmd.AddCommand(downCmd)
}
