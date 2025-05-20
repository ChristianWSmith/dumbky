package utils

import (
	"dumbky/internal/log"
	"encoding/json"

	"github.com/yosssi/gohtml"
)

func SmartFormat(text string) string {
	content := []byte(text)
	var rawJson json.RawMessage
	err := json.Unmarshal(content, &rawJson)
	if err != nil {
		log.Debug("Failed to Unmarshal data to json")
		return gohtml.Format(text)
	}
	log.Debug("Successfully Unmarshalled data to json")
	var prettyJson []byte
	prettyJson, err = json.MarshalIndent(&rawJson, "", "  ")
	if err != nil {
		log.Debug("Failed to MarshalIndent data to json")
		return text
	}
	return string(prettyJson)
}
