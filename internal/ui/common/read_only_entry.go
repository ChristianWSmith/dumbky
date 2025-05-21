package common

import (
	"dumbky/internal/log"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewReadOnlyEntry(text string) (*fyne.Container, *widget.Label, binding.String) {
	labelBind := binding.NewString()
	label := widget.NewLabel("")
	label.Bind(labelBind)
	err := labelBind.Set(text)
	if err != nil {
		log.Error(err)
	}
	th := label.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	box := canvas.NewRectangle(th.Color(theme.ColorNameInputBackground, v))
	box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = th.Size(theme.SizeNameInputBorder)
	border.StrokeColor = th.Color(theme.ColorNameInputBorder, v)
	border.CornerRadius = th.Size(theme.SizeNameInputRadius)

	stack := container.NewStack(box, border, label)

	return stack, label, labelBind
}
