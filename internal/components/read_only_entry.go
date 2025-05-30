package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ReadOnlyEntry struct {
	widget.BaseWidget
	label *widget.Label
	entry *widget.Entry
}

func NewReadOnlyEntry(text string) *ReadOnlyEntry {
	entry := widget.NewEntry()
	label := widget.NewLabel(text)
	entry.Disable()
	r := &ReadOnlyEntry{
		label: label,
		entry: entry,
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

	stack := container.NewStack(r.entry, r.label)

	return &readOnlyEntryRenderer{
		stack: stack,
		label: r.label,
		entry: r.entry,
	}
}

func (r *ReadOnlyEntry) SetText(text string) {
	r.label.SetText(text)
}

func (r *ReadOnlyEntry) GetText() string {
	return r.label.Text
}

type readOnlyEntryRenderer struct {
	stack *fyne.Container
	label *widget.Label
	entry *widget.Entry
}

func (r *readOnlyEntryRenderer) Layout(size fyne.Size) {
	r.stack.Resize(size)
}

func (r *readOnlyEntryRenderer) MinSize() fyne.Size {
	return fyne.NewSize(
		fyne.Max(r.label.MinSize().Width, r.entry.MinSize().Width),
		fyne.Max(r.label.MinSize().Height, r.entry.MinSize().Height))
}

func (r *readOnlyEntryRenderer) Refresh() {
	r.stack.Refresh()
}

func (r *readOnlyEntryRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.stack}
}

func (r *readOnlyEntryRenderer) Destroy() {
}
