package request

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type view struct {
	ui fyne.CanvasObject
}

func newView(queryParamsUI, pathParamsUI, headersUI, bodyUI fyne.CanvasObject) *view {

	queryParamsTab := container.NewTabItem(constants.UI_LABEL_QUERY_PARAMETERS, queryParamsUI)
	pathParamsTab := container.NewTabItem(constants.UI_LABEL_PATH_PARAMETERS, pathParamsUI)
	headersTab := container.NewTabItem(constants.UI_LABEL_HEADERS, headersUI)
	bodyTab := container.NewTabItem(constants.UI_LABEL_BODY, bodyUI)

	tabs := container.NewAppTabs(queryParamsTab, pathParamsTab, headersTab, bodyTab)
	ui := container.NewBorder(nil, nil, nil, nil, tabs)

	return &view{
		ui: ui,
	}
}

func (v *view) canvasObject() fyne.CanvasObject {
	return v.ui
}
