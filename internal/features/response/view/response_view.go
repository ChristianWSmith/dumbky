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
	status *components.FancyLabel
	time   *components.FancyLabel
	body   *components.FancyLabel

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

func styleLabel(label *components.FancyLabel) {
	label.SetSelectable(true)
	label.SetWrapping(fyne.TextWrapWord)
	label.SetTextStyle(fyne.TextStyle{Monospace: true})
	label.Refresh()
}

func NewView() View {
	statusEntry := components.NewFancyLabel(constants.UI_PLACEHOLDER_RESPONSE_STATUS)
	timeEntry := components.NewFancyLabel(constants.UI_PLACEHOLDER_RESPONSE_TIME)
	bodyEntry := components.NewFancyLabel(constants.UI_PLACEHOLDER_RESPONSE_BODY)

	styleLabel(statusEntry)
	styleLabel(timeEntry)
	styleLabel(bodyEntry)

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
