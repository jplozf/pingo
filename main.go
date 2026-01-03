package main

// ****************************************************************************
// IMPORTS
// ****************************************************************************
import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
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
var statusLight *canvas.Circle
var statusMutex sync.Mutex

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
	wid1 := NewPingWidget("Username:", "Enter name...", func(val string) {
		if val == "" {
			showStatus("Error: Username cannot be empty!")
		} else {
			showStatus("Username saved successfully!")
		}
	})
	wid2 := NewPingWidget("Username:", "Enter name...", func(val string) {
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
	rightContent := container.NewVBox(wid1, wid2)

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
	statusLight = canvas.NewCircle(color.NRGBA{R: 76, G: 175, B: 80, A: 255})
	statusLight.StrokeColor = color.NRGBA{R: 0, G: 0, B: 0, A: 180} // Dark charcoal/black outline
	statusLight.StrokeWidth = 1
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(16, 16)) // This forces the "slot" to be 16x16
	lightContainer := container.NewStack(rect, statusLight)
	centeredLight := container.NewCenter(lightContainer)
	barContent := container.NewHBox(
		centeredLight,
		statusLabel,
	)

	line := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 30}) // Very light grey
	line.SetMinSize(fyne.NewSize(0, 1))

	return container.NewVBox(
		line,
		barContent,
	)
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
	statusMutex.Lock()
	now := time.Now()
	lastMessageTime = now
	statusMutex.Unlock()

	fyne.Do(func() {
		statusLabel.SetText(message)

		// Change light to Yellow if not Ready
		if message != StatusDefaultMessage {
			statusLight.FillColor = color.NRGBA{R: 255, G: 235, B: 59, A: 255} // Yellow
		} else {
			statusLight.FillColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255} // Green
		}

		statusLight.Refresh()
		statusLabel.Refresh()
	})

	// Reset timer logic...
	go func() {
		time.Sleep(StatusTimeout * time.Second)
		statusMutex.Lock()
		isLast := lastMessageTime.Equal(now)
		statusMutex.Unlock()

		if isLast {
			fyne.Do(func() {
				statusLabel.SetText(StatusDefaultMessage)
				statusLight.FillColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255} // Back to Green
				statusLight.Refresh()
				statusLabel.Refresh()
			})
		}
	}()
}
