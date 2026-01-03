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
	"fyne.io/fyne/v2/dialog"
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

// ****************************************************************************
// main()
// ****************************************************************************
func main() {
	a = app.NewWithID(AppID)
	a.Settings().SetTheme(&MyCustomTheme{}) // Apply globally
	title := fmt.Sprintf("%s - v%s", AppTitle, GetDisplayVersion())
	w = a.NewWindow(title)
	w.SetIcon(theme.ComputerIcon())

	var width, height, splitOffset float64
	newConfig, err := loadConfig()
	if err != nil {
		width = 400
		height = 300
		splitOffset = 0.33
		showStatus("No config found")
	} else {
		width = float64(newConfig.WindowWidth)
		height = float64(newConfig.WindowHeight)
		splitOffset = float64(newConfig.SplitOffset)
	}
	w.Resize(fyne.NewSize(float32(width), float32(height)))

	// Save geometry when the window is closed
	w.SetOnClosed(func() {
		currSize := w.Content().Size()
		newConfig := AppConfig{
			WindowWidth:  currSize.Width,
			WindowHeight: currSize.Height,
			SplitOffset:  split.Offset,
		}
		saveConfig(newConfig)
	})

	// Bind F3 to Exit
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyF3 {
			// fmt.Println("F3 pressed: Shutting down...")
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
	// We can bind our F3 logic here too
	// quitItem := fyne.NewMenuItem("Quit (F3)", func() { w.Close() })
	// fileMenu := fyne.NewMenu("File", newItem, fyne.NewMenuItemSeparator(), quitItem)
	fileMenu := fyne.NewMenu("File", newItem)

	// Help Menu
	aboutItem := fyne.NewMenuItem("About", func() {
		dialog.ShowInformation("Pingo", "Version: "+Version+"\nAuthor: JPL", w)
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
