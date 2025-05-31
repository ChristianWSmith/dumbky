package view

import (
	"dumbky/internal/components"
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/utils"

	"dumbky/internal/features/workspace/common"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type viewImpl struct {
	ui                  *fyne.Container
	exchangeTabs        *container.DocTabs
	documentViewMap     map[*container.TabItem]common.DocumentId
	requestNameEntry    *widget.Entry
	collectionNameLabel *components.FancyLabel
	addButton           *widget.Button
	saveButton          *widget.Button
	importButton        *widget.Button
	exportButton        *widget.Button
}

type Bindables struct {
	RequestName, CollectionName features.StringBindable
}

type View interface {
	GetBindables() Bindables
	CanvasObject() fyne.CanvasObject
	GetDocumentId(tabItem *container.TabItem) common.DocumentId
	GetSelectedDocumentId() common.DocumentId
	SetSelectedDocumentText(requestName string)
	DestroyDocumentTab(tabItem *container.TabItem)
	SetDocumentClosedHandler(handler func(*container.TabItem))
	SetDocumentSelectedHandler(handler func(*container.TabItem))
	RefreshTabs()
	SelectDocumentTabOnCondition(handler func(common.DocumentId) bool) bool
	AddDocumentTab(id common.DocumentId, requestName string, ui fyne.CanvasObject)
	SetAddHandler(handler func())
	SetSaveHandler(handler func())
	SetImportHandler(handler func())
	SetExportHandler(handler func())
}

func NewView() View {

	exchangeTabs := container.NewDocTabs()
	exchangeTabs.SetTabLocation(container.TabLocationLeading)
	requestNameEntry := widget.NewEntry()

	addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)

	saveButton := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), nil)

	importButton := widget.NewButtonWithIcon("", theme.LoginIcon(), nil)
	exportButton := widget.NewButtonWithIcon("", theme.LogoutIcon(), nil)

	collectionNameLabel := components.NewFancyLabel("")
	collectionNameLabel.SetTextStyle(fyne.TextStyle{
		Bold: true,
	})
	collectionNameLabel.SetBackgroundColor(theme.Color(theme.ColorNameFocus))
	collectionNameLabel.SetBorderColor(theme.Color(theme.ColorNamePrimary))
	collectionNameLabel.Refresh()

	controlsLeft := container.NewHBox(addButton, collectionNameLabel)
	controlsRight := container.NewHBox(saveButton, importButton, exportButton)

	workspaceHeaderUI := container.NewBorder(nil, nil, controlsLeft, controlsRight, requestNameEntry)

	ui := container.NewBorder(workspaceHeaderUI, nil, nil, nil, exchangeTabs)
	return &viewImpl{
		ui:                  ui,
		exchangeTabs:        exchangeTabs,
		documentViewMap:     make(map[*container.TabItem]common.DocumentId),
		requestNameEntry:    requestNameEntry,
		collectionNameLabel: collectionNameLabel,
		addButton:           addButton,
		saveButton:          saveButton,
		importButton:        importButton,
		exportButton:        exportButton,
	}
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{
		RequestName:    v.requestNameEntry,
		CollectionName: v.collectionNameLabel,
	}
}

func (v *viewImpl) GetDocumentId(tabItem *container.TabItem) common.DocumentId {
	return v.documentViewMap[tabItem]
}

func (v *viewImpl) GetSelectedDocumentId() common.DocumentId {
	selectedTab := v.exchangeTabs.Selected()
	return v.documentViewMap[selectedTab]
}

func (v *viewImpl) SetSelectedDocumentText(requestName string) {
	v.exchangeTabs.Selected().Text = formatTabText(requestName)
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

func (v *viewImpl) AddDocumentTab(id common.DocumentId, requestName string, ui fyne.CanvasObject) {
	exchangeViewTab := container.NewTabItem(formatTabText(requestName), ui)
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

func (v *viewImpl) SetImportHandler(handler func()) {
	v.importButton.OnTapped = handler
}

func (v *viewImpl) SetExportHandler(handler func()) {
	v.exportButton.OnTapped = handler
}

func formatTabText(requestName string) string {
	return utils.TruncateStringWithEllipsis(requestName, constants.UI_DOCUMENT_TAB_MAX_LENGTH)
}
