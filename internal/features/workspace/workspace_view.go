package workspace

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type view struct {
	ui               *fyne.Container
	exchangeTabs     *container.DocTabs
	documentViewMap  map[*container.TabItem]documentId
	requestNameEntry *widget.Entry
	addButton        *widget.Button
	saveButton       *widget.Button
}

func newView() *view {

	exchangeTabs := container.NewDocTabs()
	requestNameEntry := widget.NewEntry()

	addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)

	saveButton := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), nil)

	controlsLeft := container.NewHBox(addButton)
	controlsRight := container.NewHBox(saveButton)

	workspaceHeaderUI := container.NewBorder(nil, nil, controlsLeft, controlsRight, requestNameEntry)

	ui := container.NewBorder(workspaceHeaderUI, nil, nil, nil, exchangeTabs)
	return &view{
		ui:               ui,
		exchangeTabs:     exchangeTabs,
		documentViewMap:  make(map[*container.TabItem]documentId),
		requestNameEntry: requestNameEntry,
		addButton:        addButton,
		saveButton:       saveButton,
	}
}

func (v *view) canvasObject() fyne.CanvasObject {
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

func (v *view) addDocumentTab(id documentId, collectionName, requestName string, ui fyne.CanvasObject) {
	exchangeViewTab := container.NewTabItem(formatTabText(collectionName, requestName), ui)
	v.documentViewMap[exchangeViewTab] = id
	v.exchangeTabs.Append(exchangeViewTab)
	v.exchangeTabs.Select(exchangeViewTab)
}

func formatTabText(collectionName, requestName string) string {
	return fmt.Sprintf("%s / %s", collectionName, requestName)
}

func (v *view) setAddHandler(handler func()) {
	v.addButton.OnTapped = handler
}

func (v *view) setSaveHandler(handler func()) {
	v.saveButton.OnTapped = handler
}
