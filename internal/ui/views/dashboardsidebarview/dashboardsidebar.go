package dashboardsidebarview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type DashboardSidebarView struct {
	UI               *fyne.Container
	workingDirectory string
}

func ComposeDashboardSidebarView() DashboardSidebarView {

	boundStringList := binding.NewStringList()

	list := widget.NewListWithData(boundStringList,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	ui := container.NewBorder(nil, nil, nil, nil, list)
	boundStringList.Append("A")
	boundStringList.Append("B")
	boundStringList.Append("C")
	boundStringList.Append("D")
	boundStringList.Append("E")
	boundStringList.Append("F")
	boundStringList.Append("G")
	boundStringList.Append("H")
	boundStringList.Append("I")
	boundStringList.Append("J")
	boundStringList.Append("K")
	boundStringList.Append("L")
	boundStringList.Append("M")
	boundStringList.Append("N")
	boundStringList.Append("O")
	boundStringList.Append("P")
	boundStringList.Append("Q")
	boundStringList.Append("R")
	boundStringList.Append("S")
	boundStringList.Append("T")
	boundStringList.Append("U")
	boundStringList.Append("V")
	boundStringList.Append("W")
	boundStringList.Append("X")
	boundStringList.Append("Y")
	boundStringList.Append("Z")
	return DashboardSidebarView{
		UI: ui,
	}
}
