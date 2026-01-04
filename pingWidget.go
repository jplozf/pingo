package main

// ****************************************************************************
// IMPORTS
// ****************************************************************************
import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// ****************************************************************************
// TYPES
// ****************************************************************************
type PingWidget struct {
	widget.BaseWidget
	lblLost         *ColoredLabel
	lblHostname     *ColoredLabel
	lblAddress      *ColoredLabel
	lblPingValue    *ColoredLabel
	lblAverageValue *ColoredLabel
	lblMinValue     *ColoredLabel
	lblMaxValue     *ColoredLabel
	lblRequests     *ColoredLabel
	btnDelete       *SlimButton
}

type PingHeaderWidget struct {
	widget.BaseWidget
	lblLost         *ColoredLabel
	lblHostname     *ColoredLabel
	lblAddress      *ColoredLabel
	lblPingValue    *ColoredLabel
	lblAverageValue *ColoredLabel
	lblMinValue     *ColoredLabel
	lblMaxValue     *ColoredLabel
	lblRequests     *ColoredLabel
	lblDelete       *ColoredLabel
}

// ****************************************************************************
// NewPingWidget()
// ****************************************************************************
func NewPingWidget(addressIP string) *PingWidget {
	item := &PingWidget{
		lblLost:         NewColoredLabel("-", ColorLightGrey, 11, fyne.TextAlignCenter, false),
		lblHostname:     NewColoredLabel("unknown", ColorLightGrey, 11, fyne.TextAlignCenter, false),
		lblAddress:      NewColoredLabel(addressIP, ColorLightBlue, 11, fyne.TextAlignCenter, false),
		lblPingValue:    NewColoredLabel("0", ColorLightGrey, 11, fyne.TextAlignCenter, false),
		lblAverageValue: NewColoredLabel("0", ColorLightGrey, 11, fyne.TextAlignCenter, false),
		lblMinValue:     NewColoredLabel("0", ColorLightGrey, 11, fyne.TextAlignCenter, false),
		lblMaxValue:     NewColoredLabel("0", ColorLightGrey, 11, fyne.TextAlignCenter, false),
		lblRequests:     NewColoredLabel("0", ColorLightGrey, 11, fyne.TextAlignCenter, false),
		btnDelete:       NewSlimButton("Delete", func() { fmt.Printf("Hello") }),
	}

	item.ExtendBaseWidget(item) // Critical for Fyne to recognize it as a widget
	return item
}

// ****************************************************************************
// CreateRenderer()
// ****************************************************************************
func (i *PingWidget) CreateRenderer() fyne.WidgetRenderer {
	// We use a container to handle the layout of our internal components
	content := container.NewGridWithRows(1, i.lblHostname, i.lblAddress, layout.NewSpacer(), i.lblLost, layout.NewSpacer(), i.lblPingValue, i.lblAverageValue, i.lblMinValue, i.lblMaxValue, i.lblRequests, i.btnDelete)

	return widget.NewSimpleRenderer(content)
}

// ****************************************************************************
// PingHeaderWidget()
// ****************************************************************************
func NewPingHeaderWidget() *PingHeaderWidget {
	item := &PingHeaderWidget{
		lblLost:         NewColoredLabel("Lost", ColorDarkGrey, 11, fyne.TextAlignCenter, true),
		lblHostname:     NewColoredLabel("Hostname", ColorDarkGrey, 11, fyne.TextAlignCenter, true),
		lblAddress:      NewColoredLabel("Address", ColorDarkGrey, 11, fyne.TextAlignCenter, true),
		lblPingValue:    NewColoredLabel("Ping", ColorDarkGrey, 11, fyne.TextAlignCenter, true),
		lblAverageValue: NewColoredLabel("Average", ColorDarkGrey, 11, fyne.TextAlignCenter, true),
		lblMinValue:     NewColoredLabel("Min", ColorDarkGrey, 11, fyne.TextAlignCenter, true),
		lblMaxValue:     NewColoredLabel("Max", ColorDarkGrey, 11, fyne.TextAlignCenter, true),
		lblRequests:     NewColoredLabel("Requests", ColorDarkGrey, 11, fyne.TextAlignCenter, true),
		lblDelete:       NewColoredLabel("Action", ColorDarkGrey, 11, fyne.TextAlignCenter, true),
	}

	item.ExtendBaseWidget(item) // Critical for Fyne to recognize it as a widget
	return item
}

// ****************************************************************************
// CreateRenderer()
// ****************************************************************************
func (i *PingHeaderWidget) CreateRenderer() fyne.WidgetRenderer {
	// We use a container to handle the layout of our internal components
	content := container.NewGridWithRows(1, i.lblHostname, i.lblAddress, layout.NewSpacer(), i.lblLost, layout.NewSpacer(), i.lblPingValue, i.lblAverageValue, i.lblMinValue, i.lblMaxValue, i.lblRequests, i.lblDelete)

	return widget.NewSimpleRenderer(content)
}
