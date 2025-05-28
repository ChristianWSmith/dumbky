package response

import (
	"dumbky/internal/components"
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type view struct {
	status *components.ReadOnlyEntry
	time   *components.ReadOnlyEntry
	body   *components.ReadOnlyEntry

	ui          *fyne.Container
	loadingBar  *widget.ProgressBarInfinite
	statusStack *fyne.Container
}

func styleEntry(entry *components.ReadOnlyEntry) {
	entry.SetSelectable(true)
	entry.SetWrapping(fyne.TextWrapWord)
	entry.SetTextStyle(fyne.TextStyle{Monospace: true})
}

func newView() *view {
	statusEntry := components.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_STATUS)
	timeEntry := components.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_TIME)
	bodyEntry := components.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_BODY)

	styleEntry(statusEntry)
	styleEntry(timeEntry)
	styleEntry(bodyEntry)

	loadingBar := widget.NewProgressBarInfinite()
	loadingBar.Hide()

	statusStack := container.NewVBox(loadingBar, statusEntry)
	info := container.NewVBox(statusStack, timeEntry)
	scroll := components.NewScrollInterceptorWrapper(bodyEntry)

	ui := container.NewBorder(info, nil, nil, nil, scroll)

	return &view{
		status:      statusEntry,
		time:        timeEntry,
		body:        bodyEntry,
		ui:          ui,
		loadingBar:  loadingBar,
		statusStack: statusStack,
	}
}

func (v *view) canvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *view) setLoading(loading bool) {
	if loading {
		v.status.Hide()
		v.loadingBar.Start()
		v.loadingBar.Show()
	} else {
		v.loadingBar.Stop()
		v.loadingBar.Hide()
		v.status.Show()
	}
	v.statusStack.Refresh()
}
