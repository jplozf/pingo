package main

// ****************************************************************************
// IMPORTS
// ****************************************************************************
import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ****************************************************************************
// GLOBALS
// ****************************************************************************
var a fyne.App
var w fyne.Window
var split *container.Split
var statusLabel = widget.NewLabel(StatusDefaultMessage)
var lastMessageTime time.Time
var statusChan = make(chan string, 10) // Buffer of 10 messages
var settings AppSettings

// ****************************************************************************
// main()
// ****************************************************************************
func main() {
	a = app.NewWithID(AppID)
	title := fmt.Sprintf("%s - v%s", AppTitle, GetDisplayVersion())
	w = a.NewWindow(title)
	w.SetIcon(theme.ComputerIcon())

	var width, height, splitOffset float64
	var err error
	settings, err = loadSettings()
	if err != nil {
		width = 400
		height = 300
		splitOffset = 0.33
		settings.ThemePreference = "Light"
		applyTheme(a, settings.ThemePreference)
		showStatus("No settings found")
	} else {
		width = float64(settings.WindowWidth)
		height = float64(settings.WindowHeight)
		splitOffset = float64(settings.SplitOffset)
		applyTheme(a, settings.ThemePreference)
	}
	w.Resize(fyne.NewSize(float32(width), float32(height)))

	// Save geometry when the window is closed
	w.SetOnClosed(func() {
		currSize := w.Content().Size()
		settings = AppSettings{
			WindowWidth:     currSize.Width,
			WindowHeight:    currSize.Height,
			SplitOffset:     split.Offset,
			ThemePreference: settings.ThemePreference,
		}
		saveSettings(settings)
	})

	// Bind F3 to Exit
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyF3 {
			w.Close() // window.Close() triggers OnClosed; myApp.Quit() stops the whole app
		}
	})

	// Usage
	customField := NewPingWidget("Username:", "Enter name...", func(val string) {
		if val == "" {
			showStatus("Error: Username cannot be empty!")
		} else {
			showStatus("Username saved successfully!")
		}
	})

	// Left Panel (e.g., a list or navigation)
	leftContent := container.NewVBox(
		widget.NewLabel("Navigation"),
		widget.NewButton("Setting A", func() {}),
		widget.NewButton("Setting B", func() {}),
	)

	// Right Panel (e.g., your main form)
	rightContent := container.NewVBox(customField)

	// Create the Split Container
	split = container.NewHSplit(leftContent, rightContent)
	split.Offset = splitOffset

	// Assemble Layout
	statusBar := createStatusBar()
	mainLayout := container.NewBorder(nil, statusBar, nil, nil, split)
	w.SetContent(mainLayout)

	// Setup Menu
	createMainMenu(w)
	// Start the status manager
	startStatusManager()
	// Run update check in the background
	go checkForUpdates(w)
	// And the show must go on
	w.ShowAndRun()
}

// ****************************************************************************
// createMainMenu()
// ****************************************************************************
func createMainMenu(w fyne.Window) {
	// File Menu
	newItem := fyne.NewMenuItem("New", func() { fmt.Println("Menu: New") })
	settingsItem := fyne.NewMenuItem("Settings", func() {
		showSettingsDialog(w, a, &settings)
	})

	fileMenu := fyne.NewMenu("File", newItem, settingsItem)

	// Help Menu
	aboutItem := fyne.NewMenuItem("About", func() {
		showAboutDialog(w)
	})
	helpMenu := fyne.NewMenu("Help", aboutItem)

	// Set the Main Menu
	mainMenu := fyne.NewMainMenu(fileMenu, helpMenu)
	w.SetMainMenu(mainMenu)
}

// ****************************************************************************
// createStatusBar()
// ****************************************************************************
func createStatusBar() fyne.CanvasObject {
	statusLabel.Alignment = fyne.TextAlignLeading
	// We wrap it in a container to give it a slight background or border look if desired
	return container.NewHBox(statusLabel)
}

// ****************************************************************************
// startStatusManager()
// ****************************************************************************
func startStatusManager() {
	go func() {
		for msg := range statusChan {
			// Set the new message
			fyne.Do(func() {
				statusLabel.SetText(msg)
				statusLabel.Refresh()
			})

			// Wait for the timeout
			time.Sleep(time.Duration(StatusTimeout) * time.Second)

			// Only reset to "Ready" if there isn't a newer message waiting
			if len(statusChan) == 0 {
				fyne.Do(func() {
					statusLabel.SetText(StatusDefaultMessage)
					statusLabel.Refresh()
				})
			}
		}
	}()
}

// ****************************************************************************
// showStatus()
// ****************************************************************************
func showStatus(message string) {
	// Send message to the channel without blocking
	select {
	case statusChan <- message:
	default:
		// If channel is full, we skip or handle it
	}
}
