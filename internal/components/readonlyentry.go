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

type ReadOnlyEntry struct {
	widget.BaseWidget
	label *widget.Label
}

func NewReadOnlyEntry(text string) *ReadOnlyEntry {
	r := &ReadOnlyEntry{
		label: widget.NewLabel(text),
	}
	r.ExtendBaseWidget(r)
	return r
}

func (r *ReadOnlyEntry) Bind(data binding.String) {
	r.label.Bind(data)
}

func (r *ReadOnlyEntry) SetSelectable(selectable bool) {
	r.label.Selectable = selectable
}

func (r *ReadOnlyEntry) SetWrapping(wrapping fyne.TextWrap) {
	r.label.Wrapping = wrapping
}

func (r *ReadOnlyEntry) SetTextStyle(textStyle fyne.TextStyle) {
	r.label.TextStyle = textStyle
}

func (r *ReadOnlyEntry) CreateRenderer() fyne.WidgetRenderer {
	th := r.label.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	box := canvas.NewRectangle(th.Color(theme.ColorNameInputBackground, v))
	box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = th.Size(theme.SizeNameInputBorder)
	border.StrokeColor = th.Color(theme.ColorNameInputBorder, v)
	border.CornerRadius = th.Size(theme.SizeNameInputRadius)

	stack := container.NewStack(box, border, r.label)

	return &readOnlyEntryRenderer{
		box:    box,
		border: border,
		stack:  stack,
		label:  r.label,
	}
}

func (r *ReadOnlyEntry) SetText(text string) {
	r.label.SetText(text)
}

func (r *ReadOnlyEntry) GetText() string {
	return r.label.Text
}

type readOnlyEntryRenderer struct {
	box    *canvas.Rectangle
	border *canvas.Rectangle
	stack  *fyne.Container
	label  *widget.Label
}

func (r *readOnlyEntryRenderer) Layout(size fyne.Size) {
	r.stack.Resize(size)
}

func (r *readOnlyEntryRenderer) MinSize() fyne.Size {
	return r.label.MinSize()
}

func (r *readOnlyEntryRenderer) Refresh() {
	r.stack.Refresh()
}

func (r *readOnlyEntryRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.stack}
}

func (r *readOnlyEntryRenderer) Destroy() {
}
