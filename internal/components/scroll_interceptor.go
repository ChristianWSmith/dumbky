package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ScrollInterceptor struct {
	widget.BaseWidget
	target *container.Scroll
}

func NewScrollInterceptorWrapper(canvasObject fyne.CanvasObject) fyne.CanvasObject {
	scroll := container.NewVScroll(canvasObject)
	return container.NewStack(scroll, NewScrollInterceptor(scroll))
}

func NewScrollInterceptor(target *container.Scroll) *ScrollInterceptor {
	s := &ScrollInterceptor{target: target}
	s.ExtendBaseWidget(s)
	return s
}

func (s *ScrollInterceptor) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.Transparent)
	return widget.NewSimpleRenderer(rect)
}

func (s *ScrollInterceptor) Scrolled(ev *fyne.ScrollEvent) {
	s.target.Scrolled(ev)
}
