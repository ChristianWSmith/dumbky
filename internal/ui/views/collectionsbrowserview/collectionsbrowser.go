package collectionsbrowserview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type CollectionsBrowserView struct {
	UI *fyne.Container
}

func ComposeCollectionsBrowserView() CollectionsBrowserView {
	label := widget.NewLabel("collections browser")

	ui := container.NewBorder(nil, nil, nil, nil, label)
	return CollectionsBrowserView{
		UI: ui,
	}
}
