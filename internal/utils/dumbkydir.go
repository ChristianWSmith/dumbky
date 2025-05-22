package utils

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"os"
	"path/filepath"
)

func GetDumbkyDir() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		log.Error(err)
		return "", err
	}

	dumbkyDir := filepath.Join(home, "Documents", constants.APP_NAME)

	err = os.MkdirAll(dumbkyDir, os.ModePerm)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return dumbkyDir, err
}
