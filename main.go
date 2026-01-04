package main

// ****************************************************************************
// IMPORTS
// ****************************************************************************
import (
	"fmt"
	"image/color"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
		applyTheme(a, w, settings.ThemePreference)
		showStatus("No settings found")
	} else {
		width = float64(settings.WindowWidth)
		height = float64(settings.WindowHeight)
		splitOffset = float64(settings.SplitOffset)
		applyTheme(a, w, settings.ThemePreference)
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
	wid1 := NewPingWidget("192.168.1.254")
	wid2 := NewPingWidget("8.8.8.8")

	// Left Panel (e.g., a list or navigation)
	leftContent := container.NewVBox(
		widget.NewLabel("Navigation"),
		widget.NewButton("Setting A", func() {}),
		widget.NewButton("Setting B", func() {}),
	)

	// Right Panel (e.g., your main form)
	rightContent := container.NewVBox(NewPingHeaderWidget(), wid1, wid2, layout.NewSpacer())

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
	statusLight = canvas.NewCircle(ColorGreen)
	statusLight.StrokeColor = ColorBlack
	statusLight.StrokeWidth = 1
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(16, 16)) // This forces the "slot" to be 16x16
	lightContainer := container.NewStack(rect, statusLight)
	centeredLight := container.NewCenter(lightContainer)

	versionLabel := widget.NewLabel(Version)
	versionLabel.TextStyle = fyne.TextStyle{Italic: true} // Make it look distinct

	barContent := container.NewHBox(
		centeredLight,
		statusLabel,
		layout.NewSpacer(), // PUSHES everything apart
		versionLabel,       // Stays on the RIGHT
	)

	line := canvas.NewRectangle(ColorLightGrey)
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
			statusLight.FillColor = ColorYellow
		} else {
			statusLight.FillColor = ColorGreen
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
				statusLight.FillColor = ColorGreen
				statusLight.Refresh()
				statusLabel.Refresh()
			})
		}
	}()
}

// ****************************************************************************
// GetPingTime()
// ****************************************************************************
func GetPingTime(target string) (string, error) {
	delimiter := settings.PingDelimiter
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", target)
	} else {
		cmd = exec.Command("ping", "-c", "1", target)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Look for the "time=" string in the output
	output := string(out)
	if strings.Contains(output, delimiter) {
		// Simple logic to extract the part after "time="
		parts := strings.Split(output, delimiter)
		if len(parts) > 1 {
			// Get everything after the delimiter, then grab the first word (the number)
			afterDelimiter := strings.TrimSpace(parts[1])
			timeValue := strings.Split(afterDelimiter, " ")
			return timeValue[0], nil // e.g. "14.2"
		}
	}

	return "unknown", nil
}
