package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyCustomTheme struct {
	fyne.Theme
}

// Override the Size function
func (m MyCustomTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 11 // Default is usually 14 or 15
	}
	if name == theme.SizeNameCaptionText {
		return 9 // Smaller text for captions/status bar
	}
	if name == theme.SizeNameInlineIcon {
		return 16 // You can also shrink icons if they feel too big
	}

	return theme.DefaultTheme().Size(name)
}

// Override the Color function
func (m MyCustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	// Background of the window
	case theme.ColorNameBackground:
		return color.NRGBA{R: 245, G: 246, B: 247, A: 255}
	// Main text color
	case theme.ColorNameForeground:
		return color.NRGBA{R: 45, G: 50, B: 55, A: 255}
	// Color of buttons and active elements
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 33, G: 150, B: 243, A: 255}
	// Background of Entry/Input fields
	case theme.ColorNameInputBackground:
		return color.White
	// Color of separators and borders
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 220, G: 220, B: 220, A: 255}
	}

	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

// Required interface methods (using defaults)
func (m MyCustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
func (m MyCustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
