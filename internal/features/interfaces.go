package features

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type Controller interface {
	CanvasObject() fyne.CanvasObject
}

type StringBindable interface {
	Bind(binding.String)
}

type BoolBindable interface {
	Bind(binding.Bool)
}

type StringListBindable interface {
	Bind(binding.StringList)
}
