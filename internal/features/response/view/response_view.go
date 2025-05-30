package view

import (
	"dumbky/internal/components"
	"dumbky/internal/constants"
	"dumbky/internal/features"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type viewImpl struct {
	status *components.ReadOnlyEntry
	time   *components.ReadOnlyEntry
	body   *components.ReadOnlyEntry

	ui          *fyne.Container
	loadingBar  *widget.ProgressBarInfinite
	statusStack *fyne.Container
}

type Bindables struct {
	Status, Time, Body features.StringBindable
}

type View interface {
	GetBindables() Bindables
	CanvasObject() fyne.CanvasObject
	SetLoading(loading bool)
}

func styleEntry(entry *components.ReadOnlyEntry) {
	entry.SetSelectable(true)
	entry.SetWrapping(fyne.TextWrapWord)
	entry.SetTextStyle(fyne.TextStyle{Monospace: true})
}

func NewView() View {
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

	return &viewImpl{
		status:      statusEntry,
		time:        timeEntry,
		body:        bodyEntry,
		ui:          ui,
		loadingBar:  loadingBar,
		statusStack: statusStack,
	}
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{
		Status: v.status,
		Time:   v.time,
		Body:   v.body,
	}
}

func (v *viewImpl) SetLoading(loading bool) {
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
