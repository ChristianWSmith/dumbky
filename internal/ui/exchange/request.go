package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/ui/common"

	"fyne.io/fyne/v2/container"
)

type RequestView struct {
	UI *container.AppTabs
	Params common.KeyValueEditorView
	Headers common.KeyValueEditorView
	Body RequestBodyView
}

func ComposeRequestView() RequestView {
	paramsView := common.ComposeKeyValueEditorView()
	headersView := common.ComposeKeyValueEditorView()
	bodyView := ComposeRequestBodyView()

	paramsTab := container.NewTabItem(constants.UI_LABEL_PARAMETERS, paramsView.UI)
	headersTab := container.NewTabItem(constants.UI_LABEL_HEADERS, headersView.UI)
	bodyTab := container.NewTabItem(constants.UI_LABEL_BODY, bodyView.UI)
	
	ui := container.NewAppTabs(paramsTab, headersTab, bodyTab)
	
	return RequestView {
		ui,
		paramsView,
		headersView,
		bodyView,
	}
}