package main

// ****************************************************************************
// IMPORTS
// ****************************************************************************
import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ****************************************************************************
// CONSTANTS
// ****************************************************************************
const (
	AppTitle         = "Pingo"
	AppFolderName    = ".pingo"
	SettingsFileName = "config.ini"
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
// main()
// ****************************************************************************
func main() {
	// 1. Initialize with a unique ID for persistent storage
	a := app.NewWithID("fr.ozf.pingo")
	// Use the Version variable in the title bar
	title := fmt.Sprintf("%s - v%s", AppTitle, Version)
	w := a.NewWindow(title)
	w.SetIcon(theme.ComputerIcon())

	// 2. Retrieve saved geometry (fall back to 400x300 if not found)
	width := a.Preferences().FloatWithFallback("win_width", 400)
	height := a.Preferences().FloatWithFallback("win_height", 300)
	w.Resize(fyne.NewSize(float32(width), float32(height)))

	// 3. Save geometry when the window is closed
	w.SetOnClosed(func() {
		currSize := w.Content().Size()
		a.Preferences().SetFloat("win_width", float64(currSize.Width))
		a.Preferences().SetFloat("win_height", float64(currSize.Height))
	})

	output := canvas.NewText(time.Now().Format(time.TimeOnly), color.NRGBA{R: 0xff, G: 0xff, A: 0xff})
	output.TextStyle.Monospace = true
	output.TextSize = 32
	w.SetContent(output)

	// Usage
	customField := NewPingWidget("Username:", "Enter name...", func(val string) {
		fmt.Println("Submitted value:", val)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Fyne Custom Widget"),
		customField,
	))

	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			fyne.Do(func() {
				output.Text = time.Now().Format(time.TimeOnly)
				output.Refresh()
			})
		}
	}()
	w.ShowAndRun()
}

// ****************************************************************************
// getAppFolderPath()
// ****************************************************************************
func getAppFolderPath(folderName string) (string, error) {
	// 1. Get the user's home directory (e.g., /home/user or C:\Users\user)
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// 2. Define the full path
	appPath := filepath.Join(home, folderName)

	// 3. Create the folder with 0755 permissions (rwxr-xr-x)
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
