package main

// ****************************************************************************
// IMPORTS
// ****************************************************************************
import (
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// ****************************************************************************
// TYPES
// ****************************************************************************
type AboutItem struct {
	Label string
	Value string
}

// ****************************************************************************
// GLOBALS
// ****************************************************************************
var AboutItems = []AboutItem{
	{"Author", Author},
	{"URL", AppURL},
	{"Version", Version},
	{"Built with", "Fyne & Go"},
}

// ****************************************************************************
// showAboutDialog()
// ****************************************************************************
func showAboutDialog(w fyne.Window) {
	form := container.New(layout.NewFormLayout())

	for _, item := range AboutItems {
		label := widget.NewLabel(item.Label)
		label.TextStyle = fyne.TextStyle{Bold: true}
		form.Add(label)

		if strings.HasPrefix(item.Value, "http") {
			u, _ := url.Parse(item.Value)
			form.Add(widget.NewHyperlink(item.Value, u))
		} else {
			form.Add(widget.NewLabel(item.Value))
		}
	}

	content := container.NewPadded(form)
	dialog.NewCustom("About "+AppTitle, "Close", content, w).Show()
}
