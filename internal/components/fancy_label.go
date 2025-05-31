package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FancyLabel struct {
	widget.BaseWidget
	label       *widget.Label
	bgColor     color.Color
	borderColor color.Color
}

func NewFancyLabel(text string) *FancyLabel {
	fl := &FancyLabel{}

	fl.ExtendBaseWidget(fl)
	fl.label = widget.NewLabel(text)

	return fl
}

func (fl *FancyLabel) CreateRenderer() fyne.WidgetRenderer {
	th := fl.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	box := canvas.NewRectangle(th.Color(theme.ColorNameInputBackground, v))
	box.CornerRadius = th.Size(theme.SizeNameInputRadius)

	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = th.Size(theme.SizeNameInputBorder)
	border.StrokeColor = th.Color(theme.ColorNameInputBorder, v)
	border.CornerRadius = th.Size(theme.SizeNameInputRadius)

	stack := container.NewStack(box, border, fl.label)

	return &fancyLabelRenderer{
		fl:     fl,
		border: border,
		box:    box,
		stack:  stack,
	}
}

func (fl *FancyLabel) Bind(data binding.String) {
	fl.label.Bind(data)
}

func (fl *FancyLabel) SetBackgroundColor(bgColor color.Color) {
	fl.bgColor = bgColor
}

func (fl *FancyLabel) SetBorderColor(borderColor color.Color) {
	fl.borderColor = borderColor
}

func (fl *FancyLabel) SetTextStyle(textStyle fyne.TextStyle) {
	fl.label.TextStyle = textStyle
}

func (fl *FancyLabel) SetSelectable(selectable bool) {
	fl.label.Selectable = selectable
}

func (fl *FancyLabel) SetWrapping(wrapping fyne.TextWrap) {
	fl.label.Wrapping = wrapping
}

type fancyLabelRenderer struct {
	fl     *FancyLabel
	border *canvas.Rectangle
	box    *canvas.Rectangle
	stack  *fyne.Container
}

func (r *fancyLabelRenderer) Layout(size fyne.Size) {
	r.stack.Resize(size)
}

func (r *fancyLabelRenderer) MinSize() fyne.Size {
	return r.stack.MinSize()
}

func (r *fancyLabelRenderer) Refresh() {
	r.box.FillColor = r.fl.bgColor
	r.border.StrokeColor = r.fl.borderColor
	r.stack.Refresh()
}

func (r *fancyLabelRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *fancyLabelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.stack}
}

func (r *fancyLabelRenderer) Destroy() {
}
