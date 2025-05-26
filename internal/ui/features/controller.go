package features

import "fyne.io/fyne/v2"

type Controller interface {
	GetUI() *fyne.Container
}
