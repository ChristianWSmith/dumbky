package view

import (
	"dumbky/internal/features"
	"fmt"

	"dumbky/internal/features/workspace/common"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type viewImpl struct {
	ui               *fyne.Container
	exchangeTabs     *container.DocTabs
	documentViewMap  map[*container.TabItem]common.DocumentId
	requestNameEntry *widget.Entry
	addButton        *widget.Button
	saveButton       *widget.Button
}

type Bindables struct {
	RequestName features.StringBindable
}

type View interface {
	GetBindables() Bindables
	CanvasObject() fyne.CanvasObject
	GetDocumentId(tabItem *container.TabItem) common.DocumentId
	GetSelectedDocumentId() common.DocumentId
	SetSelectedDocumentText(collectionName, requestName string)
	DestroyDocumentTab(tabItem *container.TabItem)
	SetDocumentClosedHandler(handler func(*container.TabItem))
	SetDocumentSelectedHandler(handler func(*container.TabItem))
	RefreshTabs()
	SelectDocumentTabOnCondition(handler func(common.DocumentId) bool) bool
	AddDocumentTab(id common.DocumentId, collectionName, requestName string, ui fyne.CanvasObject)
	SetAddHandler(handler func())
	SetSaveHandler(handler func())
}

func NewView() View {

	exchangeTabs := container.NewDocTabs()
	requestNameEntry := widget.NewEntry()

	addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)

	saveButton := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), nil)

	controlsLeft := container.NewHBox(addButton)
	controlsRight := container.NewHBox(saveButton)

	workspaceHeaderUI := container.NewBorder(nil, nil, controlsLeft, controlsRight, requestNameEntry)

	ui := container.NewBorder(workspaceHeaderUI, nil, nil, nil, exchangeTabs)
	return &viewImpl{
		ui:               ui,
		exchangeTabs:     exchangeTabs,
		documentViewMap:  make(map[*container.TabItem]common.DocumentId),
		requestNameEntry: requestNameEntry,
		addButton:        addButton,
		saveButton:       saveButton,
	}
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{
		RequestName: v.requestNameEntry,
	}
}

func (v *viewImpl) GetDocumentId(tabItem *container.TabItem) common.DocumentId {
	return v.documentViewMap[tabItem]
}

func (v *viewImpl) GetSelectedDocumentId() common.DocumentId {
	selectedTab := v.exchangeTabs.Selected()
	return v.documentViewMap[selectedTab]
}

func (v *viewImpl) SetSelectedDocumentText(collectionName, requestName string) {
	v.exchangeTabs.Selected().Text = formatTabText(collectionName, requestName)
}

func (v *viewImpl) DestroyDocumentTab(tabItem *container.TabItem) {
	delete(v.documentViewMap, tabItem)
}

func (v *viewImpl) SetDocumentClosedHandler(handler func(*container.TabItem)) {
	v.exchangeTabs.OnClosed = handler
}

func (v *viewImpl) SetDocumentSelectedHandler(handler func(*container.TabItem)) {
	v.exchangeTabs.OnSelected = handler
}

func (v *viewImpl) RefreshTabs() {
	v.exchangeTabs.Refresh()
}

func (v *viewImpl) SelectDocumentTabOnCondition(handler func(common.DocumentId) bool) bool {
	for tabItem, tabId := range v.documentViewMap {
		if handler(tabId) {
			v.exchangeTabs.Select(tabItem)
			return true
		}
	}
	return false
}

func (v *viewImpl) AddDocumentTab(id common.DocumentId, collectionName, requestName string, ui fyne.CanvasObject) {
	exchangeViewTab := container.NewTabItem(formatTabText(collectionName, requestName), ui)
	v.documentViewMap[exchangeViewTab] = id
	v.exchangeTabs.Append(exchangeViewTab)
	v.exchangeTabs.Select(exchangeViewTab)
}

func (v *viewImpl) SetAddHandler(handler func()) {
	v.addButton.OnTapped = handler
}

func (v *viewImpl) SetSaveHandler(handler func()) {
	v.saveButton.OnTapped = handler
}

func formatTabText(collectionName, requestName string) string {
	return fmt.Sprintf("%s / %s", collectionName, requestName)
}
