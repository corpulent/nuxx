package common

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func GetToken() interface{} {
	cliToken := viper.Get("token")

	if cliToken == "" {
		fmt.Println("Oops, looks like you don't have an access token. Please run nuxx init to get one.")
		os.Exit(1)
	}

	return cliToken
}
