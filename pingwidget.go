package main

// ****************************************************************************
// IMPORTS
// ****************************************************************************
import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ****************************************************************************
// TYPES
// ****************************************************************************
type PingWidget struct {
	widget.BaseWidget
	Label  *widget.Label
	Entry  *widget.Entry
	Button *widget.Button
}

// ****************************************************************************
// NewPingWidget()
// ****************************************************************************
func NewPingWidget(labelPath, placeholder string, onAction func(string)) *PingWidget {
	item := &PingWidget{
		Label: widget.NewLabel(labelPath),
		Entry: widget.NewEntry(),
	}
	item.Entry.SetPlaceHolder(placeholder)
	item.Button = widget.NewButton("Submit", func() {
		onAction(item.Entry.Text)
	})

	item.ExtendBaseWidget(item) // Critical for Fyne to recognize it as a widget
	return item
}

// ****************************************************************************
// CreateRenderer()
// ****************************************************************************
func (i *PingWidget) CreateRenderer() fyne.WidgetRenderer {
	// We use a container to handle the layout of our internal components
	content := container.NewBorder(nil, nil, i.Label, i.Button, i.Entry)

	return widget.NewSimpleRenderer(content)
}
