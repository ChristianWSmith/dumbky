package utils

import (
	"dumbky/internal/constants"
	"strings"

	"github.com/sethvargo/go-diceware/diceware"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SillyName() string {
	list, _ := diceware.Generate(constants.UI_DICEWARE_COUNT)
	title := []string{}
	for _, item := range list {
		title = append(title, cases.Title(language.English, cases.Compact).String(item))
	}
	return strings.Join(title, "")
}
