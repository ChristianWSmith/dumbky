package utils

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"os"
	"path/filepath"
)

func getAppDataDir() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		log.Error(err)
		return "", err
	}

	appDataDir := filepath.Join(configDir, constants.APP_NAME)

	err = os.MkdirAll(appDataDir, os.ModePerm)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return appDataDir, err
}

func GetDBFilePath() (string, error) {
	appDataDir, err := getAppDataDir()
	if err != nil {
		log.Error(err)
		return "", err
	}
	dbFilePath := filepath.Join(appDataDir, constants.DB_FILE_NAME)
	return dbFilePath, nil
}
