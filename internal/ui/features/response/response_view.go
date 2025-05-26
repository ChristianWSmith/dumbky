package response

import (
	"dumbky/internal/constants"
	"dumbky/internal/ui/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type view struct {
	statusLabel *widget.Label
	timeLabel   *widget.Label
	bodyLabel   *widget.Label

	ui              *fyne.Container
	statusContainer *fyne.Container
	loadingBar      *widget.ProgressBarInfinite
	statusStack     *fyne.Container
}

func styleLabel(label *widget.Label) {
	label.Selectable = true
	label.Wrapping = fyne.TextWrapWord
	label.TextStyle.Monospace = true
}

func newView() *view {
	statusContainer, statusLabel := components.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_STATUS)
	timeContainer, timeLabel := components.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_TIME)
	bodyContainer, bodyLabel := components.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_BODY)

	styleLabel(statusLabel)
	styleLabel(timeLabel)
	styleLabel(bodyLabel)

	loadingBar := widget.NewProgressBarInfinite()
	loadingBar.Hide()

	statusStack := container.NewVBox(loadingBar, statusContainer)
	info := container.NewVBox(statusStack, timeContainer)
	ui := container.NewBorder(info, nil, nil, nil, container.NewVScroll(bodyContainer))

	return &view{
		statusLabel:     statusLabel,
		timeLabel:       timeLabel,
		bodyLabel:       bodyLabel,
		ui:              ui,
		statusContainer: statusContainer,
		loadingBar:      loadingBar,
		statusStack:     statusStack,
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}

func (v *view) setLoading(loading bool) {
	if loading {
		v.statusContainer.Hide()
		v.loadingBar.Start()
		v.loadingBar.Show()
	} else {
		v.loadingBar.Stop()
		v.loadingBar.Hide()
		v.statusContainer.Show()
	}
	v.statusStack.Refresh()
}
