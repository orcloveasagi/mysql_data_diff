package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeNav() fyne.CanvasObject {
	app := fyne.CurrentApp()

	tree := &widget.Tree{
		CreateNode: func(branch bool) (o fyne.CanvasObject) {
			return widget.NewLabel("Menu")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			m, ok := Menus[uid]
			if !ok {
				fyne.LogError("Missing Menu: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(m.Title)
		},
		ChildUIDs: func(uid widget.TreeNodeID) (c []widget.TreeNodeID) {
			menu := Menus[uid]
			if menu == nil {
				return nil
			}
			if menu.Children == nil {
				return nil
			}
			return menu.Children
		},
		IsBranch: isBranch,
		OnSelected: func(uid string) {
			if menu, ok := Menus[uid]; ok {
				app.Preferences().SetString("selectedMenu", uid)
				if menu.Onselect != nil {
					forms := menu.Onselect(&NavCtx{uid: uid})
					updateForm(forms)
				}
			}
		},
	}

	currentPref := app.Preferences().StringWithFallback("selectedMenu", CreateMenuUID)
	tree.Select(currentPref)

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			app.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
		}),
		widget.NewButton("Light", func() {
			app.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantLight})
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}

func isBranch(uid string) bool {
	if menu, ok := Menus[uid]; ok {
		return menu.Children != nil
	}
	return false
}
