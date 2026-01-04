package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ****************************************************************************
// TYPES
// ****************************************************************************
type AppSettings struct {
	WindowWidth     float32 `json:"window_width"`
	WindowHeight    float32 `json:"window_height"`
	SplitOffset     float64 `json:"split_offset"`
	ThemePreference string  `json:"theme_preference"` // "Light" or "Dark"
	PingDelimiter   string  `json:"ping_delimiter"`
}

// ****************************************************************************
// saveSettings()
// ****************************************************************************
func saveSettings(settings AppSettings) error {
	path, _ := getAppFolderPath(AppFolderName)
	// Create folder if missing
	os.MkdirAll(filepath.Dir(path), 0755)
	// Convert struct to JSON bytes
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(path, SettingsFileName), data, 0644)
}

// ****************************************************************************
// loadSettings()
// ****************************************************************************
func loadSettings() (AppSettings, error) {
	path, _ := getAppFolderPath(AppFolderName)
	var settings AppSettings
	data, err := os.ReadFile(filepath.Join(path, SettingsFileName))
	if err != nil {
		return settings, err // Return empty config if file doesn't exist
	}
	err = json.Unmarshal(data, &settings)
	return settings, err
}

// ****************************************************************************
// getAppFolderPath()
// ****************************************************************************
func getAppFolderPath(folderName string) (string, error) {
	// Get the user's home directory (e.g., /home/user or C:\Users\user)
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// Define the full path
	appPath := filepath.Join(home, folderName)
	// Create the folder with 0755 permissions (rwxr-xr-x)
	// If it exists, MkdirAll returns nil (no error)
	err = os.MkdirAll(appPath, 0755)
	if err != nil {
		return "", err
	}
	return appPath, nil
}

// ****************************************************************************
// showSettingsDialog()
// ****************************************************************************
func showSettingsDialog(parentWin fyne.Window, myApp fyne.App, settings *AppSettings) {
	// 1. Existing Theme Selection
	themeSelect := widget.NewSelect([]string{"Light", "Dark"}, func(value string) {
		applyTheme(myApp, parentWin, value)
		settings.ThemePreference = value
		saveSettings(*settings)
	})
	themeSelect.SetSelected(settings.ThemePreference)

	// 2. New Ping Delimiter Entry
	pingEntry := widget.NewEntry()
	pingEntry.SetText(settings.PingDelimiter)
	pingEntry.PlaceHolder = "e.g., time= or temps="

	// This function saves the setting as the user types
	pingEntry.OnChanged = func(value string) {
		settings.PingDelimiter = value
		saveSettings(*settings)
	}

	// 3. Assemble the Content
	content := container.NewVBox(
		widget.NewLabel("Preferred Theme:"),
		themeSelect,
		widget.NewSeparator(), // Adds a nice line between sections
		widget.NewLabel("Ping 'Time' Delimiter (OS Specific):"),
		pingEntry,
		widget.NewLabelWithStyle("(English: 'time=', French: 'temps=')",
			fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
	)

	d := dialog.NewCustom("Settings", "Close", content, parentWin)
	// We increase the height slightly to fit the new fields
	d.Resize(fyne.NewSize(350, 250))
	d.Show()
}

// ****************************************************************************
// applyTheme()
// ****************************************************************************
func applyTheme(myApp fyne.App, myWin fyne.Window, preference string) {
	switch preference {
	case "Dark":
		myApp.Settings().SetTheme(theme.DarkTheme())
	case "Light":
		myApp.Settings().SetTheme(theme.LightTheme())
	default:
		myApp.Settings().SetTheme(theme.DefaultTheme())
	}
	myWin.Content().Refresh()
}
