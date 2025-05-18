package exchange

import (
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
	statusEntry.Disable()
	statusEntry.SetPlaceHolder("<response status>")
	timeEntry := widget.NewEntry()
	timeEntry.Disable()
	timeEntry.SetPlaceHolder("<response time>")
	bodyEntry := widget.NewMultiLineEntry()
	bodyEntry.Disable()
	bodyEntry.SetPlaceHolder("<response body>")

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