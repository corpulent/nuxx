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

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve deployed service or job logs.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(projectConfFile); os.IsNotExist(err) {
			d := color.New(color.FgCyan, color.Bold).SprintFunc()
			fmt.Printf(`Project configuration file does not exist! Create it by running: %s`, d("nuxx compose"))
			os.Exit(1)
		}

		cliToken := common.GetToken()

		p := fmt.Println
		file, _ := ioutil.ReadFile(projectConfFile)
		json.Unmarshal([]byte(file), &configStructureMap)
		existingData = gojsonq.New().File(projectConfFile)
		projectName := existingData.Find("project_name")
		releaseName, _ := cmd.Flags().GetString("release-name")
		releaseKind, _ := cmd.Flags().GetString("release-kind")
		stringCyan := color.New(color.FgCyan, color.Bold).SprintFunc()

		if releaseName == "" {
			var jsonStr = fmt.Sprintf(`{"project_name": "%s"}`, projectName)
			releasesEndpoint := fmt.Sprintf(`%s%s`, releasesEndpoint, cliToken)
			releasesResp := pkg.ReleasesRequest(releasesEndpoint, jsonStr)

			if len(releasesResp.Resp.Jobs) > 0 {
				p(``)
				p(`Jobs executions in this project:`)
				p(``)

				latestString := "latest"
				activeString := "running"

				for _, v := range releasesResp.Resp.Jobs {
					latest := v.Latest
					active := v.Active

					if latest == 0 {
						latestString = ""
					}

					if active == 0 {
						activeString = ""
					}

					p(v.RELEASE_NAME, latestString, activeString)
				}

				p(``)
				p(`Run`, stringCyan("$ nuxx status -k job -n EXECUTION_NAME"), `to check the status.`)
				p(``)
			}

			if len(releasesResp.Resp.Services) > 0 {
				p(``)
				p(`Service releases in this project:`)
				p(``)

				for _, v := range releasesResp.Resp.Services {
					p(v.RELEASE_NAME)
				}

				p(``)
				p(`Run`, stringCyan("$ nuxx logs -k service -n RELEASE_NAME"), `to check the status.`)
				p(``)
			}
		}

		if releaseName != "" && releaseKind != "" {
			logsJsonStr := fmt.Sprintf(`{"project_name": "%s", "service_name": "%s"}`, projectName, releaseName)

			if releaseKind == "job" {
				logsJsonStr = fmt.Sprintf(`{"project_name": "%s", "execution_id": "%s"}`, projectName, releaseName)
			}

			if releaseKind == "service" {
				logsJsonStr = fmt.Sprintf(`{"project_name": "%s", "service_name": "%s"}`, projectName, releaseName)
			}

			logsPathEndpoint := fmt.Sprintf(`%s%s`, logsPathEndpoint, cliToken)
			resp := pkg.LogsRequest(logsPathEndpoint, logsJsonStr)
			printAllLogsResponses(resp, releaseName, releaseKind)
		}
	},
}

func printAllLogsResponses(resp *pkg.LogResponse, releaseName string, releaseKind string) {
	respResp := resp.Resp
	p := fmt.Println
	stringCyan := color.New(color.FgCyan, color.Bold).SprintFunc()

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

	if len(respResp.Logs) > 0 {
		p(``)
		logsFoundString := fmt.Sprintf(`Logs for %s %s...`, releaseKind, stringCyan(releaseName))
		p(logsFoundString)
		p(``)
		for _, v := range respResp.Logs {
			p(v)
		}
		p(``)
	} else {
		noLogsFoundString := fmt.Sprintf(`No logs where found for %s %s...`, releaseKind, stringCyan(releaseName))
		p(``)
		p(noLogsFoundString)
		p(``)
	}
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringP("release-name", "n", "", "Release name")
	logsCmd.Flags().StringP("release-kind", "k", "", "Release kind, job or service")
}
