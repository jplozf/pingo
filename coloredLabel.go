package main

// ****************************************************************************
// IMPORTS
// ****************************************************************************
import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ****************************************************************************
// TYPES
// ****************************************************************************
type ColoredLabel struct {
	widget.BaseWidget
	text      string
	color     color.Color
	size      float32
	alignment fyne.TextAlign
	Bold      bool
	Padding   fyne.Size
}

type coloredLabelRenderer struct {
	item    *ColoredLabel
	bg      *canvas.Rectangle
	text    *canvas.Text
	objects []fyne.CanvasObject
}

// ****************************************************************************
// NewColoredLabel()
// ****************************************************************************
func NewColoredLabel(text string, bgColor color.Color, size float32, alignment fyne.TextAlign, bold bool) *ColoredLabel {
	item := &ColoredLabel{
		text:      text,
		color:     bgColor,
		size:      size,
		alignment: alignment,
		Bold:      bold,
		Padding:   fyne.NewSize(1, 1),
	}

	item.ExtendBaseWidget(item) // Critical for Fyne to recognize it as a widget
	return item
}

// ****************************************************************************
// CreateRenderer()
// ****************************************************************************
func (i *ColoredLabel) CreateRenderer() fyne.WidgetRenderer {
	txt := canvas.NewText(i.text, GetContrastColor(i.color))
	txt.TextSize = i.size
	txt.TextStyle = fyne.TextStyle{Bold: i.Bold}

	bg := canvas.NewRectangle(i.color)
	content := container.NewStack(bg, container.NewCenter(txt))

	return &coloredLabelRenderer{
		item:    i,
		bg:      bg,
		text:    txt,
		objects: []fyne.CanvasObject{content},
	}
}

// ****************************************************************************
// MinSize()
// ****************************************************************************
func (i *ColoredLabel) MinSize() fyne.Size {
	// 1. Get the size the text actually needs
	// We create a temporary text object to measure it
	tempText := canvas.NewText(i.text, ColorBlack)
	tempText.TextSize = i.size
	textSize := tempText.MinSize()

	// 2. Add your custom padding to that size
	// For a "small" label, try something like 4px horizontal, 2px vertical
	return fyne.NewSize(
		textSize.Width+i.Padding.Width,
		textSize.Height+i.Padding.Height,
	)
}

// ****************************************************************************
// Refresh()
// ****************************************************************************
func (r *coloredLabelRenderer) Refresh() {
	// 1. Get the current background color from the widget
	bgColor := r.item.color

	// 2. Update the background rectangle
	r.bg.FillColor = bgColor

	// 3. FORCE re-calculation of contrast based on the CURRENT background
	r.text.Color = GetContrastColor(bgColor)

	// 4. Update text and style
	r.text.Text = r.item.text
	r.text.TextStyle.Bold = r.item.Bold

	r.bg.Refresh()
	r.text.Refresh()
}

// ****************************************************************************
// GetContrastColor()
// ****************************************************************************
// GetContrastColor returns white or black depending on the background brightness
func GetContrastColor(bg color.Color) color.Color {
	r, g, b, _ := bg.RGBA()

	// These constants (0.299, 0.587, 0.114) represent how the human eye
	// perceives the brightness of Red, Green, and Blue.
	// We divide by 65535 because Go's RGBA() returns 16-bit values.
	lum := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535

	if lum < 0.5 {
		return color.White // Background is dark, use white text
	}
	return color.Black // Background is light, use black text
}

// ****************************************************************************
// Layout()
// ****************************************************************************
// 1. Mandatory: How to arrange the internal objects
func (r *coloredLabelRenderer) Layout(size fyne.Size) {
	// This tells the stack (background + text) to fill the whole widget area
	r.objects[0].Resize(size)
}

// ****************************************************************************
// MinSize()
// ****************************************************************************
// 2. Mandatory: How big the widget wants to be
func (r *coloredLabelRenderer) MinSize() fyne.Size {
	return r.item.MinSize() // This calls the MinSize we wrote for the widget earlier
}

// ****************************************************************************
// Objects()
// ****************************************************************************
// 3. Mandatory: Return the list of objects to draw
func (r *coloredLabelRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// ****************************************************************************
// Destroy()
// ****************************************************************************
// 4. Mandatory: Cleanup when the widget is removed (The one you were missing!)
func (r *coloredLabelRenderer) Destroy() {
	// Usually left empty unless you have high-memory resources to free
}

// ****************************************************************************
// SetText()
// ****************************************************************************
func (i *ColoredLabel) SetText(newText string) {
	i.text = newText
	i.Refresh() // This triggers r.Refresh() in the renderer
}

// ****************************************************************************
// SetColor()
// ****************************************************************************
func (i *ColoredLabel) SetColor(newColor color.Color) *ColoredLabel {
	i.color = newColor
	i.Refresh() // This tells the renderer: "The data changed, redraw now!"
	return i
}

// ****************************************************************************
// SetBold()
// ****************************************************************************
func (i *ColoredLabel) SetBold(b bool) *ColoredLabel {
	i.Bold = b
	i.Refresh()
	return i
}
