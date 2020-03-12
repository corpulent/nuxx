package cmd

import (
	"bufio"
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

var (
	respData           = &pkg.Response{}
	serviceData        = &pkg.Service{}
	configStructureMap map[string]interface{}
	configStructure    = `{
		"project_name" : "",
		"services": {}
	}`
)

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Create a nuxx.json configuration file.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildProjectConfig()
	},
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func buildProjectConfig() error {
	// check if a configuration file already exists,  then unmarshall it
	// if not, then use the default data from configStructure
	if fileExists(projectConfFile) {
		fmt.Printf("A configuration file already exists!  Do you want to continue? y/N: ")
		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()

		if char == 'Y' || char == 'y' {
			fmt.Println("Moving on...")
		} else {
			fmt.Println("Understood! Aborting...")
			os.Exit(1)
		}

		file, _ := ioutil.ReadFile(projectConfFile)
		json.Unmarshal([]byte(file), &configStructureMap)
	} else {
		json.Unmarshal([]byte(configStructure), &configStructureMap)
		jsonBytes, _ := json.MarshalIndent(configStructureMap, "", "  ")
		_ = ioutil.WriteFile(projectConfFile, jsonBytes, 0644)
	}

	existingData = gojsonq.New().File(projectConfFile)

	// keep sending data to the server until the server deems the configuration complete
	for {
		jsonBytes, _ := json.MarshalIndent(configStructureMap, "", "  ")
		jsonString := string(jsonBytes)
		done, resp := ask(jsonString)
		stringCyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		respDataErr := json.Unmarshal([]byte(resp), &respData)
		common.CheckError(respDataErr)
		askFor := respData.Resp.ASK_FOR

		if done == "true" {
			fmt.Println(``)
			fmt.Println(`Your project looks complete!  Run `, stringCyan("$ nuxx up"), ` to deploy.`)
			fmt.Println(``)
			break
		}

		if askFor == "project_name" {
			projectName, err := common.PromptString("Name your project")
			common.CheckError(err)
			configStructureMap["project_name"] = projectName
			_, err = os.OpenFile(projectConfFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
			common.CheckError(err)
		}

		if askFor == "services_jobs" {
			fmt.Println(``)
			fmt.Println(`Your project configuration is created, but you need to add some services or jobs.`)
			fmt.Println(``)
			fmt.Println(`Visit `, stringCyan("https://nuxx.io/getting-started"), ` to get started.`)
			fmt.Println(``)
			break
		}
	}

	// save the file
	jsonBytes, _ := json.MarshalIndent(configStructureMap, "", "  ")
	writeFileErr := ioutil.WriteFile(projectConfFile, jsonBytes, 0644)
	common.CheckError(writeFileErr)

	return nil
}

func ask(payload string) (string, string) {
	resp := pkg.PostRequest(composeEndpoint, payload)
	_ = json.Unmarshal([]byte(resp), &respData)
	complete := respData.Resp.Complete
	return complete, resp
}

func init() {
	rootCmd.AddCommand(composeCmd)
}
