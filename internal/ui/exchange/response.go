package exchange

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ResponseView struct {
	UI *fyne.Container
	Status binding.String
	Time binding.String
	Body binding.String
}

func ComposeResponseView() ResponseView {
	statusBind := binding.NewString()
	timeBind := binding.NewString()
	bodyBind := binding.NewString()

	statusEntry := widget.NewEntry()
	statusEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_RESPONSE_STATUS)
	statusEntry.TextStyle.Monospace = true
	timeEntry := widget.NewEntry()
	timeEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_RESPONSE_TIME)
	timeEntry.TextStyle.Monospace = true
	bodyEntry := widget.NewEntry()
	bodyEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_RESPONSE_BODY)
	bodyEntry.TextStyle.Monospace = true

	statusEntry.Bind(statusBind)
	timeEntry.Bind(timeBind)
	bodyEntry.Bind(bodyBind)

	info := container.NewVBox(statusEntry, timeEntry)
	ui := container.NewBorder(info, nil, nil, nil, bodyEntry)

	return ResponseView {
		ui,
		statusBind,
		timeBind,
		bodyBind,
	}
}