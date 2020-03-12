package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/corpulent/nuxx/pkg"
	"github.com/corpulent/nuxx/pkg/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/thedevsaddam/gojsonq"
)

var (
	signedUrlRespData = &pkg.SignedUrlResponse{}
	zipDestination    string
	zipPath           string
	payload           string
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Deploy an application or run a job.",
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

		assetsPath, _ := cmd.Flags().GetString("assets-path")
		uploadPathEndpoint := fmt.Sprintf(`%s%s`, uploadPathEndpoint, cliToken)
		upPathEndpoint := fmt.Sprintf(`%s%s`, upPathEndpoint, cliToken)

		if assetsPath != "" {
			upPathEndpoint = fmt.Sprintf(`%s?assets=true`, upPathEndpoint)
			zipPath = path.Join(assetsPath)
			zipDestination = fmt.Sprintf(`%v.zip`, projectName)
			pkg.RecursiveZip(zipPath, zipDestination)

			fmt.Println(``)
			fmt.Println(`File package`, zipDestination, `created...`)

			payload = fmt.Sprintf(`{"project_name": "%s", "archive_name": "%s"}`, projectName, zipDestination)
			resp := pkg.PostRequest(uploadPathEndpoint, payload)
			_ = json.Unmarshal([]byte(resp), &signedUrlRespData)
			uploadURL := signedUrlRespData.Resp.PRESIGNED_URL

			fmt.Println("Package being uploaded... please wait.")

			uploadResp := pkg.UploadRequest(zipDestination, uploadURL)

			if uploadResp.StatusCode == 200 {
				fmt.Println("Package upload complete!")
				fmt.Println(``)
			} else {
				fmt.Println("Something went wrong during upload...")
				fmt.Println(``)
			}
		}

		jsonBytes, _ := json.MarshalIndent(configStructureMap, "", "  ")
		jsonString := string(jsonBytes)
		fmt.Println("Deploying...")
		fmt.Println(``)
		upResp := upRequest(upPathEndpoint, jsonString)
		stringCyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		stringGreen := color.New(color.FgGreen, color.Bold).SprintFunc()
		stringBold := color.New(color.Bold).SprintFunc()

		for releaseName, v := range upResp {
			releaseType := "service"
			if releaseName[:3] == "job" {
				releaseType = "job"
			}

			statusString := fmt.Sprintf("$ nuxx status -k %s -n %s", releaseType, releaseName)
			logsString := fmt.Sprintf("$ nuxx logs -k %s -n %s", releaseType, releaseName)

			fmt.Println(strings.Title(releaseType), stringBold(releaseName), stringGreen(strings.ToUpper(v.COMMAND_RESPONSE.Status)))
			fmt.Println(``)
			fmt.Println(`Check status`, stringCyan(statusString))
			fmt.Println(`View the logs`, stringCyan(logsString))
			fmt.Println(``)
			fmt.Println("If you are experiencing any issues with our cli, please email us at", stringCyan("support@nuxx.io"), ".")
			fmt.Println(``)
		}
	},
}

func upRequest(urlEndpoint string, payload string) map[string]pkg.UpRelease {
	respData := &pkg.UpResponse{}
	resp := pkg.PostRequest(urlEndpoint, payload)
	_ = json.Unmarshal([]byte(resp), &respData)

	return respData.Resp
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().String("assets-path", "", "Provide an archive assets path as part of your deployment.")
}
