package features

import "fyne.io/fyne/v2"

type Controller interface {
	CanvasObject() fyne.CanvasObject
}
