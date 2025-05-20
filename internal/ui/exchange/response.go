package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/ui/common"

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
}

func styleResponseLabel(label *widget.Label) {
	label.Selectable = true
	label.Wrapping = fyne.TextWrapWord
	label.TextStyle.Monospace = true
}

func ComposeResponseView() ResponseView {
	statusBind := binding.NewString()
	timeBind := binding.NewString()
	bodyBind := binding.NewString()

	statusEntry, statusLabel, statusBind := common.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_STATUS)
	timeEntry, timeLabel, timeBind := common.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_TIME)
	bodyEntry, bodyLabel, bodyBind := common.NewReadOnlyEntry(constants.UI_PLACEHOLDER_RESPONSE_BODY)

	styleResponseLabel(statusLabel)
	styleResponseLabel(timeLabel)
	styleResponseLabel(bodyLabel)

	statusBind.Set(constants.UI_PLACEHOLDER_RESPONSE_STATUS)

	info := container.NewVBox(statusEntry, timeEntry)
	ui := container.NewBorder(info, nil, nil, nil, container.NewVScroll(bodyEntry))

	return ResponseView{
		ui,
		statusBind,
		timeBind,
		bodyBind,
	}
}
