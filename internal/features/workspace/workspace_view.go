package workspace

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type view struct {
	ui              *fyne.Container
	exchangeTabs    *container.DocTabs
	documentViewMap map[*container.TabItem]documentId
}

func newView(workspaceHeaderUI *fyne.Container) *view {

	exchangeTabs := container.NewDocTabs()

	ui := container.NewBorder(workspaceHeaderUI, nil, nil, nil, exchangeTabs)
	return &view{
		ui:              ui,
		exchangeTabs:    exchangeTabs,
		documentViewMap: make(map[*container.TabItem]documentId),
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}

func (v *view) getDocumentId(tabItem *container.TabItem) documentId {
	return v.documentViewMap[tabItem]
}

func (v *view) getSelectedDocumentId() documentId {
	selectedTab := v.exchangeTabs.Selected()
	return v.documentViewMap[selectedTab]
}

func (v *view) setSelectedDocumentText(collectionName, requestName string) {
	v.exchangeTabs.Selected().Text = formatTabText(collectionName, requestName)
}

func (v *view) destroyDocumentTab(tabItem *container.TabItem) {
	delete(v.documentViewMap, tabItem)
}

func (v *view) setDocumentClosedHandler(handler func(*container.TabItem)) {
	v.exchangeTabs.OnClosed = handler
}

func (v *view) setDocumentSelectedHandler(handler func(*container.TabItem)) {
	v.exchangeTabs.OnSelected = handler
}

func (v *view) refreshTabs() {
	v.exchangeTabs.Refresh()
}

func (v *view) selectDocumentTabOnCondition(handler func(documentId) bool) bool {
	for tabItem, tabId := range v.documentViewMap {
		if handler(tabId) {
			v.exchangeTabs.Select(tabItem)
			return true
		}
	}
	return false
}

func (v *view) addDocumentTab(id documentId, collectionName, requestName string, ui *fyne.Container) {
	exchangeViewTab := container.NewTabItem(formatTabText(collectionName, requestName), ui)
	v.documentViewMap[exchangeViewTab] = id
	v.exchangeTabs.Append(exchangeViewTab)
	v.exchangeTabs.Select(exchangeViewTab)
}

func formatTabText(collectionName, requestName string) string {
	return fmt.Sprintf("%s / %s", collectionName, requestName)
}
