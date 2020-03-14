package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/corpulent/nuxx/pkg"
	"github.com/corpulent/nuxx/pkg/common"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thedevsaddam/gojsonq"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check deployment status.",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("release-name", cmd.Flags().Lookup("release-name"))
		viper.BindPFlag("release-kind", cmd.Flags().Lookup("release-kind"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		p := fmt.Println

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
				p(`Run`, stringCyan("$ nuxx status -k service -n RELEASE_NAME"), `to check the status.`)
				p(``)
			}
		}

		if releaseName != "" && releaseKind != "" {
			if releaseKind == "job" {
				jobStatusPathEndpoint := fmt.Sprintf(`%s%s`, jobStatusPathEndpoint, cliToken)
				jobStatusJsonStr := fmt.Sprintf(`{"project_name": "%s", "execution_id": "%s"}`, projectName, releaseName)
				resp := pkg.JobStatusRequest(jobStatusPathEndpoint, jobStatusJsonStr)

				printAllJobResponses(resp, releaseName)
			}

			if releaseKind == "service" {
				srvStatusPathEndpoint := fmt.Sprintf(`%s%s`, srvStatusPathEndpoint, cliToken)
				serviceStatusJsonStr := fmt.Sprintf(`{"project_name": "%s", "service_name": "%s"}`, projectName, releaseName)
				resp := pkg.ServiceStatusRequest(srvStatusPathEndpoint, serviceStatusJsonStr)

				printAllServiceResponses(resp, releaseName)
			}
		}
	},
}

func printAllServiceResponses(resp *pkg.ServiceStatus, releaseName string) {
	respResp := resp.Resp
	p := fmt.Println
	stringGreen := color.New(color.FgGreen, color.Bold).SprintFunc()
	stringRed := color.New(color.FgRed, color.Bold).SprintFunc()
	stringYellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	stringBold := color.New(color.Bold).SprintFunc()

	if respResp.Response.Status != "" {
		p(``)
		switch status := respResp.Response.Status; status {
		case "running":
			serviceUpStatusString := fmt.Sprintf(`Service %s is %s`, stringBold(releaseName), stringGreen("RUNNING"))
			p(serviceUpStatusString)
		case "deploying":
			deployingStatusString := fmt.Sprintf(`Service %s is %s`, stringBold(releaseName), stringYellow("DEPLOYING"))
			p(deployingStatusString)
		case "unhealthy":
			unhealthyStatusString := fmt.Sprintf(`Service %s is %s`, stringBold(releaseName), stringYellow("UNHEALTHY"))
			p(unhealthyStatusString)
			p("The service can still be running, but it might have failed a health check.  Check the status logs below.")
		default:
			serviceDownStatusString := fmt.Sprintf(`Service %s state: %s`, stringBold(releaseName), stringRed(status))
			p(serviceDownStatusString)
		}
	} else {
		p(``)
		emptyStatusString := fmt.Sprintf(`Service %s state unknown...`, stringBold(releaseName))
		p(emptyStatusString)
	}

	if len(respResp.STATUS_MESSAGES) > 0 {
		status_data := [][]string{}
		p(``)
		for _, v := range respResp.STATUS_MESSAGES {
			status_data = append(status_data, []string{stringBold(v.Type), v.Reason, v.Message})
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Type", "Status", "Message"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t") // pad with tabs
		table.SetNoWhiteSpace(true)
		table.AppendBulk(status_data)
		table.Render()
	}

	if len(respResp.Notices) > 0 {
		p(``)
		for _, v := range respResp.Notices {
			p(v)
		}
	}

	if len(respResp.Warnings) > 0 {
		p(``)
		for _, v := range respResp.Warnings {
			p(v)
		}
	}

	if len(respResp.Errors) > 0 {
		p(``)
		for _, v := range respResp.Errors {
			p(v)
		}
	}
	p(``)
}

func printAllJobResponses(resp *pkg.JobStatus, releaseName string) {
	respResp := resp.Resp
	p := fmt.Println
	jobStatus := respResp.JOB_STATUS
	exitCode := 0
	reason := ""
	timeLayout := time.RFC3339
	startAtString := ""
	endAtString := ""
	stringBold := color.New(color.Bold).SprintFunc()

	if jobStatus.ExitCode != 0 {
		exitCode = jobStatus.ExitCode
	}

	if jobStatus.Reason != "" {
		reason = jobStatus.Reason
	}

	if jobStatus.StartedAt != "" {
		startAtTime, startAtTimeErr := time.Parse(timeLayout, jobStatus.StartedAt)

		if startAtTimeErr == nil {
			startAtString = fmt.Sprintf(`Started at %s.`, startAtTime)
		}
	}

	if jobStatus.StartedAt != "" {
		endAtTime, endAtTimeErr := time.Parse(timeLayout, jobStatus.FinishedAt)

		if endAtTimeErr == nil {
			endAtString = fmt.Sprintf(`Finished at %s`, endAtTime)
		}
	}

	p(``)
	statusString := fmt.Sprintf(`Job %s with exit code %d.`, strings.ToLower(reason), exitCode)
	p(statusString)

	if startAtString != "" {
		p(startAtString)
	}

	if endAtString != "" {
		p(endAtString)
	}

	if len(respResp.STATUS_MESSAGES) > 0 {
		status_data := [][]string{}
		p(``)
		for _, v := range respResp.STATUS_MESSAGES {
			status_data = append(status_data, []string{stringBold(v.Type), v.Reason, v.Message})
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Type", "Status", "Message"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t") // pad with tabs
		table.SetNoWhiteSpace(true)
		table.AppendBulk(status_data)
		table.Render()
	}

	if len(respResp.Notices) > 0 {
		p(``)
		for _, v := range respResp.Notices {
			p(v)
		}
	}

	if len(respResp.Warnings) > 0 {
		p(``)
		for _, v := range respResp.Warnings {
			p(v)
		}
	}

	if len(respResp.Errors) > 0 {
		p(``)
		for _, v := range respResp.Errors {
			p(v)
		}
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringP("release-name", "n", "", "Release name")
	statusCmd.Flags().StringP("release-kind", "k", "", "Release kind, job or service")
}
