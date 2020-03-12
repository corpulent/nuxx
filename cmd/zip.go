package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/corpulent/nuxx/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/thedevsaddam/gojsonq"
)

var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "Create a zip archive of your files.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(projectConfFile); os.IsNotExist(err) {
			d := color.New(color.FgCyan, color.Bold).SprintFunc()
			fmt.Printf(`Project configuration file does not exist! Create it by running: %s`, d("nuxx compose"))
			os.Exit(1)
		}

		flagZipName, _ := cmd.Flags().GetString("zip-name")
		dir, _ := os.Getwd()
		var zipName string
		var zipPath string

		if flagZipName != "" {
			zipName = flagZipName
			zipPath = fmt.Sprintf(`%s`, path.Join(dir, zipName))
		} else {
			file, _ := ioutil.ReadFile(projectConfFile)
			json.Unmarshal([]byte(file), &configStructureMap)
			existingData = gojsonq.New().File(projectConfFile)
			projectName := existingData.Find("project_name")
			zipName = fmt.Sprintf(`%s`, projectName)
			zipPath = fmt.Sprintf(`%s`, path.Join(dir, zipName))
		}

		fmt.Println("Creating file package...")
		_ = pkg.DeleteFile(zipPath)
		filesList := pkg.GenerateFilesToZip(dir)
		pkg.ZipFiles(zipPath, filesList)
		successString := fmt.Sprintf(`File package created %s!`, zipName)
		fmt.Println(successString)
	},
}

func init() {
	rootCmd.AddCommand(zipCmd)
	zipCmd.Flags().String("zip-name", "", "Zip file name")
}
