package collectionsbrowserview

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"dumbky/internal/db"
	"dumbky/internal/global"
	"dumbky/internal/log"
)

type CollectionsBrowserView struct {
	UI              *fyne.Container
	SelectedBinding binding.String
	collectionNames []string
	listWidget      *widget.List
}

func ComposeCollectionsBrowserView() CollectionsBrowserView {
	selected := binding.NewString()
	collections := fetchCollectionNames()

	view := CollectionsBrowserView{
		SelectedBinding: selected,
		collectionNames: collections,
	}

	view.listWidget = widget.NewList(
		func() int { return len(view.collectionNames) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			menuButton := widget.NewButtonWithIcon("", nil, nil)
			menuButton.Icon = menuButton.Theme().Icon(theme.IconNameMoreVertical)
			return container.NewBorder(nil, nil, nil, menuButton, label)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			label := c.Objects[0].(*widget.Label)
			menuBtn := c.Objects[1].(*widget.Button)

			name := view.collectionNames[i]
			label.SetText(name)

			menuBtn.OnTapped = func() {
				pop := fyne.NewMenu("", // not shown
					fyne.NewMenuItem("Delete", func() {
						log.Info(fmt.Sprintf("TODO: Delete %s", name))
					}),
					fyne.NewMenuItem("Rename", func() {
						log.Info(fmt.Sprintf("TODO: Rename %s", name))
					}),
				)
				widget.ShowPopUpMenuAtPosition(pop, global.Window.Canvas(), menuBtn.Position())
			}
		},
	)

	view.listWidget.OnSelected = func(id widget.ListItemID) {
		selected.Set(view.collectionNames[id])
	}

	addBtn := widget.NewButton("Add Collection", func() {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("Collection name")
		dialog.ShowForm("New Collection", "Add", "Cancel", []*widget.FormItem{
			widget.NewFormItem("Name", entry),
		}, func(ok bool) {
			if ok && entry.Text != "" {
				err := addCollection(entry.Text)
				if err != nil {
					dialog.ShowError(err, global.Window)
					return
				}
				view.refresh()
			}
		}, global.Window)
	})

	view.UI = container.NewBorder(nil, addBtn, nil, nil, view.listWidget)
	return view
}

func fetchCollectionNames() []string {
	rows, err := db.DB.Query("SELECT name FROM collections ORDER BY created_at ASC")
	if err != nil {
		log.Error(err)
		return []string{}
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Error(err)
			continue
		}
		names = append(names, name)
	}
	return names
}

func addCollection(name string) error {
	_, err := db.DB.Exec("INSERT INTO collections (name) VALUES (?)", name)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (v *CollectionsBrowserView) refresh() {
	v.collectionNames = fetchCollectionNames()
	v.listWidget.Refresh()
}
