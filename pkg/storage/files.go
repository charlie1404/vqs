package storage

import (
	"os"
	"path"
	"syscall"

	"github.com/charlie1404/vqs/pkg/o11y/logs"
)

func createQueueFolderExists(queueDir string) error {
	err := os.RemoveAll(queueDir)
	if err != nil {
		logs.Logger.Warn().Err(err).Msg("error removing queue dir")
	}

	if err := os.MkdirAll(queueDir, 0700); err != nil {
		logs.Logger.Error().Err(err).Msg("error creating queue dir")
		return err
	}

	return nil
}

func createQueueDataFilesExists(queueDirPath string) error {
	metaFile := path.Join(queueDirPath, "meta.dat")
	dataFile := path.Join(queueDirPath, "data.dat")
	inFlightFile := path.Join(queueDirPath, "in_flight.dat")
	delayedFile := path.Join(queueDirPath, "delayed.dat")

	_ = os.Remove(metaFile)
	_ = os.Remove(dataFile)
	_ = os.Remove(inFlightFile)
	_ = os.Remove(delayedFile)

	// Meta File ===============================
	file, err := os.OpenFile(metaFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := syscall.Truncate(metaFile, int64(META_FILE_SIZE)); err != nil {
		os.Remove(metaFile)
		return err
	}

	// Data File ===============================
	file, err = os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := syscall.Truncate(dataFile, int64(DATA_BUFFER_SIZE)); err != nil {
		os.Remove(dataFile)
		return err
	}

	// // In FlightMessages Data File ===============================
	// file, err = os.OpenFile(inFlightFile, os.O_RDWR|os.O_CREATE, 0600)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// if err := syscall.Truncate(inFlightFile, INITIAL_QUEUE_FILE_SIZE); err != nil {
	// 	os.Remove(dataFile)
	// 	return err
	// }

	// // Delayed Messages File ===============================
	// file, err = os.OpenFile(delayedFile, os.O_RDWR|os.O_CREATE, 0600)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// if err := syscall.Truncate(delayedFile, INITIAL_QUEUE_FILE_SIZE); err != nil {
	// 	os.Remove(dataFile)
	// 	return err
	// }

	return nil
}
