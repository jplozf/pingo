package main

// ****************************************************************************
// IMPORTS
// ****************************************************************************
import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// ****************************************************************************
// TYPES
// ****************************************************************************
type SlimButton struct {
	widget.BaseWidget
	Text     string
	OnTapped func()
	hovering bool // Tracks the hover state
}

type slimButtonRenderer struct {
	bg      *canvas.Rectangle
	text    *canvas.Text
	objects []fyne.CanvasObject
	button  *SlimButton
}

// ****************************************************************************
// NewSlimButton()
// ****************************************************************************
func NewSlimButton(text string, tapped func()) *SlimButton {
	b := &SlimButton{
		Text:     text,
		OnTapped: tapped,
	}
	b.ExtendBaseWidget(b)
	return b
}

// ****************************************************************************
// CreateRenderer()
// ****************************************************************************
func (i *SlimButton) CreateRenderer() fyne.WidgetRenderer {
	// 1. Background (using a grey color from your colors.go eventually)
	bg := canvas.NewRectangle(color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	bg.CornerRadius = 3

	// 2. Text (forced to a smaller size)
	txt := canvas.NewText(i.Text, color.Black)
	txt.TextSize = 11 // Smaller than default 14
	txt.Alignment = fyne.TextAlignCenter

	return &slimButtonRenderer{
		bg:      bg,
		text:    txt,
		objects: []fyne.CanvasObject{bg, txt},
		button:  i,
	}
}

// ****************************************************************************
// Layout()
// ****************************************************************************
func (r *slimButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	// Center the text manually within the available height
	r.text.Resize(size)
}

// ****************************************************************************
// MinSize()
// ****************************************************************************
func (r *slimButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.text.MinSize().Width+12, r.text.MinSize().Height+4)
}

// ****************************************************************************
// Objects()
// ****************************************************************************
func (r *slimButtonRenderer) Objects() []fyne.CanvasObject { return r.objects }

// ****************************************************************************
// Destroy()
// ****************************************************************************
func (r *slimButtonRenderer) Destroy() {}

// ****************************************************************************
// Tapped()
// ****************************************************************************
func (i *SlimButton) Tapped(_ *fyne.PointEvent) {
	if i.OnTapped != nil {
		i.OnTapped()
	}
}

// ****************************************************************************
// TappedSecondary()
// ****************************************************************************
func (i *SlimButton) TappedSecondary(_ *fyne.PointEvent) {} // Required for interface

// ****************************************************************************
// MouseIn()
// ****************************************************************************
func (i *SlimButton) MouseIn(_ *desktop.MouseEvent) {
	i.hovering = true
	i.Refresh() // Triggers the Renderer's Refresh()
}

// ****************************************************************************
// MouseOut()
// ****************************************************************************
func (i *SlimButton) MouseOut() {
	i.hovering = false
	i.Refresh()
}

// ****************************************************************************
// MouseMoved()
// ****************************************************************************
func (i *SlimButton) MouseMoved(_ *desktop.MouseEvent) {} // Not needed but part of interface

// ****************************************************************************
// Refresh()
// ****************************************************************************
func (r *slimButtonRenderer) Refresh() {
	if r.button.hovering {
		// Darker grey when hovering
		r.bg.FillColor = color.NRGBA{R: 180, G: 180, B: 180, A: 255}
	} else {
		// Standard light grey
		r.bg.FillColor = color.NRGBA{R: 210, G: 210, B: 210, A: 255}
	}

	r.text.Text = r.button.Text
	r.bg.Refresh()
	r.text.Refresh()
}
