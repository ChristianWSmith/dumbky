package dashboardsidebarview

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type DashboardSidebarView struct {
	UI *fyne.Container
	//workingDirectory string
}

func ComposeDashboardSidebarView() DashboardSidebarView {
	currentPath, _ := os.Getwd()
	pathLabel := widget.NewLabel(currentPath)

	fileList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {},
	)

	var updateFileList func(string)

	updateFileList = func(path string) {
		files, err := os.ReadDir(path)
		if err != nil {
			pathLabel.SetText("Failed to read directory")
			return
		}

		items := []os.FileInfo{}
		items = append(items, nil) // For ".."

		for _, file := range files {
			info, _ := file.Info()
			items = append(items, info)
		}

		fileList.Length = func() int {
			return len(items)
		}
		fileList.UpdateItem = func(i widget.ListItemID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if i == 0 {
				label.SetText("..")
			} else {
				label.SetText(items[i].Name())
			}
		}
		fileList.OnSelected = func(i widget.ListItemID) {
			fileList.UnselectAll()
			if i == 0 {
				newPath := filepath.Dir(path)
				currentPath = newPath
				pathLabel.SetText(newPath)
				updateFileList(newPath)
			} else {
				selected := items[i]
				fullPath := filepath.Join(path, selected.Name())
				if selected.IsDir() {
					currentPath = fullPath
					pathLabel.SetText(fullPath)
					updateFileList(fullPath)
				} else {
					fmt.Println("Selected file:", fullPath)
				}
			}
		}
		fileList.Refresh()
	}

	updateFileList(currentPath)

	ui := container.NewBorder(pathLabel, nil, nil, nil, fileList)
	return DashboardSidebarView{
		UI: ui,
	}
}
