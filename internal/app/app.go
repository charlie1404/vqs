package app

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/charlie1404/vqs/internal/api"
	"github.com/charlie1404/vqs/internal/o11y/logs"
	"github.com/charlie1404/vqs/internal/o11y/metrics"
)

func clearScreen() {
	// https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
	fmt.Printf("\x1B[H\x1B[2J\x1B[3J")
}

func ensureDataFolderExist() error {
	logs.Logger.Info().Msg("checking data dir")

	var dirName string = "data"
	var folderMode os.FileMode = 0700

	info, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		logs.Logger.Info().Msg("data dir does not exist, creating one")
		if err := os.Mkdir(dirName, folderMode); err != nil {
			logs.Logger.Error().Err(err).Msg("")
			return err
		}
		logs.Logger.Info().Msg("created data dir")
		return nil
	}

	if !info.IsDir() {
		return errors.New("data dir path exists but is not a directory")
	}

	if info.Mode() != folderMode {
		if err := os.Chmod(dirName, folderMode); err != nil {
			logs.Logger.Error().Err(err).Msg("error setting data dir permissions")
			return errors.New("error setting data dir permissions")
		}
	}
	logs.Logger.Info().Msg("data dir exists")
	return nil
}

func New() {
	clearScreen()
	logs.InitLogger()

	err := ensureDataFolderExist()
	if err != nil {
		logs.Logger.Fatal().Err(err).Msg("Failed to create data folder")
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		metrics := metrics.New()
		metrics.StartServer()
		metrics.SetupInterruptListener()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		apiApp := api.New()
		apiApp.StartServer()
		apiApp.SetupInterruptListener()
		apiApp.CloseQueues()
		wg.Done()
	}()

	wg.Wait()
}
