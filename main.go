package main

// ****************************************************************************
// IMPORTS
// ****************************************************************************
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ****************************************************************************
// TYPES
// ****************************************************************************
type AppConfig struct {
	WindowWidth  float32 `json:"window_width"`
	WindowHeight float32 `json:"window_height"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
}

// ****************************************************************************
// GLOBALS
// ****************************************************************************
var a fyne.App
var statusLabel = widget.NewLabel("Ready")
var lastMessageTime time.Time
var statusChan = make(chan string, 10) // Buffer of 10 messages

// ****************************************************************************
// main()
// ****************************************************************************
func main() {
	// Initialize with a unique ID for persistent storage
	a = app.NewWithID(AppID)
	// Use the Version variable in the title bar
	title := fmt.Sprintf("%s - v%s", AppTitle, GetDisplayVersion())
	w := a.NewWindow(title)
	w.SetIcon(theme.ComputerIcon())

	// Retrieve saved geometry (fall back to 400x300 if not found)
	width := a.Preferences().FloatWithFallback("win_width", 400)
	height := a.Preferences().FloatWithFallback("win_height", 300)
	w.Resize(fyne.NewSize(float32(width), float32(height)))

	// Save geometry when the window is closed
	w.SetOnClosed(func() {
		currSize := w.Content().Size()
		a.Preferences().SetFloat("win_width", float64(currSize.Width))
		a.Preferences().SetFloat("win_height", float64(currSize.Height))
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

	// Setup Menu
	createMainMenu(w)

	// Assemble Layout
	statusBar := createStatusBar()
	mainContent := container.NewVBox(customField)

	// Border Layout: Top=nil, Bottom=statusBar, Left=nil, Right=nil, Center=mainContent
	layout := container.NewBorder(nil, statusBar, nil, nil, mainContent)

	w.SetContent(layout)

	// Start the manager
	startStatusManager()

	// Run update check in the background
	go checkForUpdates(w)

	w.ShowAndRun()
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
// saveData()
// ****************************************************************************
func saveData(content string) error {
	folder, err := getAppFolderPath(AppFolderName) // Starting with a dot makes it hidden on Linux/macOS
	if err != nil {
		return err
	}

	filePath := filepath.Join(folder, SettingsFileName)

	// WriteFile creates or overwrites the file
	return os.WriteFile(filePath, []byte(content), 0644)
}

// ****************************************************************************
// saveConfig()
// ****************************************************************************
func saveConfig(config AppConfig) error {
	path, _ := getAppFolderPath(AppFolderName)

	// Create folder if missing
	os.MkdirAll(filepath.Dir(path), 0755)

	// Convert struct to JSON bytes
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ****************************************************************************
// loadConfig()
// ****************************************************************************
func loadConfig() (AppConfig, error) {
	path, _ := getAppFolderPath(AppFolderName)
	var config AppConfig

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err // Return empty config if file doesn't exist
	}

	err = json.Unmarshal(data, &config)
	return config, err
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
