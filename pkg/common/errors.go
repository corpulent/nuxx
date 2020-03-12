package common

import (
	"log"
	"os"
)

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
