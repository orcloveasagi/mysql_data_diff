package gui

import (
	"fyne.io/fyne/v2/widget"
)

func NewEntry(data string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(data)
	return entry
}

func NewPswEntry(data string) *widget.Entry {
	entry := NewEntry(data)
	entry.Password = true
	return entry
}

func FlipForm(formA *widget.Form, formB *widget.Form) {
	for i, sourceItem := range formA.Items {
		sourceEntry := sourceItem.Widget.(*widget.Entry)
		targetEntry := formB.Items[i].Widget.(*widget.Entry)

		sourceText := sourceEntry.Text
		sourceEntry.SetText(targetEntry.Text)
		targetEntry.SetText(sourceText)
	}
}
