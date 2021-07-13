package app

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/charlie1404/vqueue/pkg/api"
)

func ensureDataFolderExist() error {
	dirName := "./data"

	log.Println("checking data dir")

	createErr := os.Mkdir(dirName, 0755)

	if createErr != nil {
		if !os.IsExist(createErr) {
			log.Printf("error creating data dir: %+v\n", createErr)
			return createErr
		}

		info, statErr := os.Stat(dirName)
		if statErr != nil {
			log.Printf("error getting stats of data dir: %+v\n", statErr)
			return statErr
		}

		if !info.IsDir() {
			log.Println("path exists but is not a directory")
			return errors.New("path exists but is not a directory")
		}

		log.Printf("data dir already exist\n")

		modeErr := os.Chmod(dirName, 0755)

		if modeErr != nil {
			log.Printf("data dir not writeable\n")
			return modeErr
		}
	}

	log.Println("created data dir")
	return nil
}

func New() {
	// https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
	fmt.Printf("\x1B[H\x1B[2J\x1B[3J")

	log.SetOutput(os.Stdout)

	err := ensureDataFolderExist()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("starting api server")

	apiApp := api.New()
	apiApp.StartServer()
}
