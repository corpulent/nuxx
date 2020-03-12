package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/corpulent/nuxx/pkg"
	"github.com/spf13/cobra"
	"github.com/thedevsaddam/gojsonq"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile               string
	baseURL               = "https://nuxx.io"
	composeEndpoint       = fmt.Sprintf(`%s/api/compose/`, baseURL)
	getTokenEndpoint      = fmt.Sprintf(`%s/api/access/`, baseURL)
	uploadPathEndpoint    = fmt.Sprintf(`%s/api/upload_path/`, baseURL)
	upPathEndpoint        = fmt.Sprintf(`%s/api/action/up/`, baseURL)
	downPathEndpoint      = fmt.Sprintf(`%s/api/action/down/`, baseURL)
	logsPathEndpoint      = fmt.Sprintf(`%s/api/logs/`, baseURL)
	srvStatusPathEndpoint = fmt.Sprintf(`%s/api/status/`, baseURL)
	jobStatusPathEndpoint = fmt.Sprintf(`%s/api/job_status/`, baseURL)
	releasesEndpoint      = fmt.Sprintf(`%s/api/releases/`, baseURL)
	projectConfFile       = "./nuxx.json"
	tokenRespData         = &pkg.TokenResponse{}
	emailPayloadMap       map[string]interface{}
	existingData          *gojsonq.JSONQ
	authToken             string
)

var rootCmd = &cobra.Command{
	Use:   "nuxx",
	Short: "",
	Long: `
	Nuxx CLI helps with fast micro-service deployments.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/nuxx.yaml)")
}

func initConfig() {
	home, err := homedir.Dir()
	configPath := fmt.Sprintf("%s/nuxx.yaml", home)

	viper.SetConfigName("nuxx.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)

	if cfgFile != "" {
		abs, err := filepath.Abs(cfgFile)
		if err != nil {
			log.Fatal("Error reading filepath: ", err.Error())
		}

		base := filepath.Base(abs)
		path := filepath.Dir(abs)

		viper.SetConfigName(strings.Split(base, ".")[0])
		viper.AddConfigPath(path)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Failed to read the global config file: ", err.Error())
			os.Exit(1)
		}
	} else {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				_, file_err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)

				if file_err != nil {
					log.Fatal(file_err)
				}

				fmt.Println("Global config file created in", home)
				viper.SetConfigFile(configPath)
				viper.Set("token", "")
				if err := viper.WriteConfig(); err != nil {
					fmt.Println(err.Error())
				}
			} else {
				fmt.Println("Config file was found but another error was produced")
			}
		}
	}

	viper.AutomaticEnv()
}
