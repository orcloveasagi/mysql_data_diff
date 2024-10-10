package gui

import (
	"db_diff/db"
	"db_diff/logic"
	"db_diff/util"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	formDb   = "db"
	formName = "name"
)

func createForm(data *db.CompareData) *fyne.Container {
	sourceForm := dbData2Form(data.Source)
	targetForm := dbData2Form(data.Target)
	commonForm := commonData2Form(data.Common)
	updateName(sourceForm, targetForm, commonForm)
	flipButton := widget.NewButton("FLIP", func() {
		FlipForm(sourceForm, targetForm)
	})
	execButton := widget.NewButton("EXEC", func() {
		onExec(sourceForm, targetForm, commonForm, data.Id)
	})

	sourceBox := container.NewVBox(widget.NewLabelWithStyle("SOURCE", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), sourceForm)
	targetBox := container.NewVBox(widget.NewLabelWithStyle("TARGET", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), targetForm)
	dbDataContainer := container.NewGridWithColumns(3, sourceBox, flipButton, targetBox)
	return container.NewBorder(dbDataContainer, nil, nil, nil,
		layout.NewSpacer(),
		container.NewVBox(widget.NewLabelWithStyle("COMMON", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			commonForm,
			widget.NewLabel(""),
			container.NewGridWithColumns(3, layout.NewSpacer(), execButton, layout.NewSpacer())),
	)
}

func updateName(source *widget.Form, target *widget.Form, common *widget.Form) {
	var name string
	for i, item := range source.Items {
		if item.Text == formDb {
			sourceEntry := item.Widget.(*widget.Entry)
			targetEntry := target.Items[i].Widget.(*widget.Entry)
			name = sourceEntry.Text + " ==> " + targetEntry.Text
			if sourceEntry.OnChanged == nil {
				sourceEntry.OnChanged = func(s string) {
					updateName(source, target, common)
				}
			}
			if targetEntry.OnChanged == nil {
				targetEntry.OnChanged = func(s string) {
					updateName(source, target, common)
				}
			}
			break
		}
	}
	for _, item := range common.Items {
		if item.Text == formName {
			item.Widget.(*widget.Entry).SetText(name)
			break
		}
	}
}

func dbData2Form(data *db.Db) *widget.Form {
	return widget.NewForm(
		widget.NewFormItem(formDb, NewEntry(data.Db)),
		widget.NewFormItem("host", NewEntry(data.Host)),
		widget.NewFormItem("port", NewEntry(data.Port)),
		widget.NewFormItem("user", NewEntry(data.User)),
		widget.NewFormItem("psw", NewPswEntry(data.Psw)),
	)
}

func form2DbData(form *widget.Form) *db.Db {
	formKV := make(map[string]string)
	for _, value := range form.Items {
		formKV[value.Text] = value.Widget.(*widget.Entry).Text
	}
	return &db.Db{
		Db:   formKV[formDb],
		Host: formKV["host"],
		Port: formKV["port"],
		User: formKV["user"],
		Psw:  formKV["psw"],
	}
}

func commonData2Form(data *db.Common) *widget.Form {
	return widget.NewForm(
		widget.NewFormItem(formName, NewEntry(data.Name)),
		widget.NewFormItem("path", NewEntry(data.Path)),
		widget.NewFormItem("ddl", NewEntry(data.Ddl)),
		widget.NewFormItem("dml", NewEntry(data.Dml)),
	)
}

func form2CommonData(form *widget.Form) *db.Common {
	formKV := make(map[string]string)
	for _, value := range form.Items {
		formKV[value.Text] = value.Widget.(*widget.Entry).Text
	}
	return &db.Common{
		Path: formKV["path"],
		Ddl:  formKV["ddl"],
		Dml:  formKV["dml"],
		Name: formKV[formName],
	}
}

func onExec(sourceForm, targetForm, commonForm *widget.Form, id int64) {
	formData := &db.CompareData{
		Source: form2DbData(sourceForm),
		Target: form2DbData(targetForm),
		Common: form2CommonData(commonForm),
		Id:     id,
	}
	dir, err := logic.DatabaseDiff(formData)
	if err != nil {
		fyne.LogError("exec error", err)
		showError(err)
		return
	}
	err = util.OpenExplorer(dir)
	if err != nil {
		fyne.LogError("open explorer", err)
	}
	if id < 0 {
		err := db.Insert(formData)
		if err != nil {
			fyne.LogError("insert compare data error", err)
		}
		refreshMenuData()
		updateNav(makeNav())
	} else {
		err := db.Update(formData)
		if err != nil {
			fyne.LogError("update  compare data error", err)
		}
		refreshMenuData()
		updateNav(makeNav())
	}
	showSuccess(dir)
}
