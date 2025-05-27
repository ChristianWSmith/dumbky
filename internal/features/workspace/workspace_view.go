package workspace

import (
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

func (v *view) setSelectedDocumentText(text string) {
	v.exchangeTabs.Selected().Text = text
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
