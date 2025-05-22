package dashboardsidebarview

import (
	"dumbky/internal/log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type DashboardSidebarView struct {
	UI                      *fyne.Container
	WorkingDirectoryBinding binding.String
	SelectedFileBinding     binding.String

	fileList *widget.List
}

func ComposeDashboardSidebarView() DashboardSidebarView {
	currentPath, _ := os.Getwd()

	workingDirectoryBind := binding.NewString()
	err := workingDirectoryBind.Set(currentPath)
	if err != nil {
		log.Error(err)
	}

	pathLabel := widget.NewLabel("")
	pathLabel.Bind(workingDirectoryBind)

	selectedFileBinding := binding.NewString()

	fileList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {},
	)

	ui := container.NewBorder(pathLabel, nil, nil, nil, fileList)

	dsv := DashboardSidebarView{
		UI:                      ui,
		WorkingDirectoryBinding: workingDirectoryBind,
		fileList:                fileList,
		SelectedFileBinding:     selectedFileBinding,
	}
	go dsv.updateFileList()

	return dsv
}

func (dsv DashboardSidebarView) updateFileList() {
	workingDirectory, err := dsv.WorkingDirectoryBinding.Get()
	if err != nil {
		log.Error(err)
		return
	}

	items, err := collectFileInfo(workingDirectory)
	if err != nil {
		log.Error(err)
		return
	}

	fyne.Do(func() {
		dsv.updateFileListWidget(items, workingDirectory)
	})
}

func (dsv DashboardSidebarView) updateFileListWidget(items []os.FileInfo, workingDirectory string) {
	dsv.fileList.Length = func() int {
		return len(items)
	}
	dsv.fileList.UpdateItem = func(i widget.ListItemID, o fyne.CanvasObject) {
		label := o.(*widget.Label)
		if i == 0 {
			label.SetText("..")
		} else {
			label.SetText(items[i].Name())
		}
	}
	dsv.fileList.OnSelected = func(i widget.ListItemID) {
		dsv.fileList.UnselectAll()
		if i == 0 {
			err := dsv.WorkingDirectoryBinding.Set(filepath.Dir(workingDirectory))
			if err != nil {
				log.Error(err)
				return
			}
			go dsv.updateFileList()
		} else {
			selected := items[i]
			fullPath := filepath.Join(workingDirectory, selected.Name())
			if selected.IsDir() {
				err := dsv.WorkingDirectoryBinding.Set(fullPath)
				if err != nil {
					log.Error(err)
					return
				}
				go dsv.updateFileList()
			} else {
				err := dsv.SelectedFileBinding.Set(fullPath)
				if err != nil {
					log.Error(err)
					return
				}
			}
		}
	}
	dsv.fileList.Refresh()
}

func collectFileInfo(workingDirectory string) ([]os.FileInfo, error) {
	items := []os.FileInfo{}

	files, err := os.ReadDir(workingDirectory)
	if err != nil {
		log.Error(err)
		return items, err
	}

	items = append(items, nil) // For ".."

	for _, file := range files {
		info, _ := file.Info()
		items = append(items, info)
	}
	return items, nil
}
