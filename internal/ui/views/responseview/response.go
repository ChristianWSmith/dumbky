package responseview

import (
	"dumbky/internal/constants"
	"dumbky/internal/ui/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ResponseView struct {
	UI            *fyne.Container
	StatusBinding binding.String
	TimeBinding   binding.String
	BodyBinding   binding.String

	statusContainer *fyne.Container
	loadingBar      *widget.ProgressBarInfinite
	statusStack     *fyne.Container
}

func styleResponseLabel(label *widget.Label) {
	label.Selectable = true
	label.Wrapping = fyne.TextWrapWord
	label.TextStyle.Monospace = true
}

func (rv ResponseView) SetLoading(loading bool) {
	if loading {
		rv.statusContainer.Hide()
		rv.loadingBar.Start()
		rv.loadingBar.Show()
		rv.statusStack.Refresh()
	} else {
		rv.loadingBar.Hide()
		rv.loadingBar.Stop()
		rv.statusContainer.Show()
		rv.statusStack.Refresh()
	}
}

func ComposeResponseView() ResponseView {
	statusContainer, statusLabel, statusBind := components.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_STATUS)
	timeContainer, timeLabel, timeBind := components.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_TIME)
	bodyContainer, bodyLabel, bodyBind := components.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_BODY)

	styleResponseLabel(statusLabel)
	styleResponseLabel(timeLabel)
	styleResponseLabel(bodyLabel)

	statusBind.Set(constants.UI_PLACEHOLDER_RESPONSE_STATUS)

	loadingBar := widget.NewProgressBarInfinite()
	loadingBar.Stop()
	loadingBar.Hide()
	statusStack := container.NewVBox(loadingBar, statusContainer)

	info := container.NewVBox(statusStack, timeContainer)
	ui := container.NewBorder(info, nil, nil, nil, container.NewVScroll(bodyContainer))

	return ResponseView{
		ui,
		statusBind,
		timeBind,
		bodyBind,
		statusContainer,
		loadingBar,
		statusStack,
	}
}
