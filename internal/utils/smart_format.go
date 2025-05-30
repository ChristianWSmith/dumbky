package utils

import (
	"encoding/json"

	"github.com/yosssi/gohtml"
)

func SmartFormat(text string) string {
	content := []byte(text)
	var rawJson json.RawMessage
	err := json.Unmarshal(content, &rawJson)
	if err != nil {
		return gohtml.Format(text)
	}
	var prettyJson []byte
	prettyJson, err = json.MarshalIndent(&rawJson, "", "  ")
	if err != nil {
		return text
	}
	return string(prettyJson)
}
