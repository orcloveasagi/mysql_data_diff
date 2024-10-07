package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var (
	MainWindow fyne.Window
	Navigation *fyne.Container
	Form       *fyne.Container
)

const appName = "Database Differ"

func Run() {
	myApp := app.NewWithID(appName)
	MainWindow = myApp.NewWindow(appName)

	// Form to add/edit data
	Navigation = container.NewStack()
	Form = container.NewStack()
	split := container.NewHSplit(Navigation, Form)
	split.Offset = 0.2
	MainWindow.SetContent(split)
	MainWindow.Resize(fyne.NewSize(1366, 768))

	updateNav(makeNav())
	MainWindow.ShowAndRun()
}

func updateNav(objects ...fyne.CanvasObject) {
	Navigation.Objects = objects
	Navigation.Refresh()
}

func updateForm(objects ...fyne.CanvasObject) {
	Form.Objects = objects
	Form.Refresh()
}
