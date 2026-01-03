package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ****************************************************************************
// TYPES
// ****************************************************************************
type AppConfig struct {
	WindowWidth  float32 `json:"window_width"`
	WindowHeight float32 `json:"window_height"`
	SplitOffset  float64 `json:"split_offset"`
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
	return os.WriteFile(filepath.Join(path, SettingsFileName), data, 0644)
}

// ****************************************************************************
// loadConfig()
// ****************************************************************************
func loadConfig() (AppConfig, error) {
	path, _ := getAppFolderPath(AppFolderName)
	var config AppConfig
	data, err := os.ReadFile(filepath.Join(path, SettingsFileName))
	if err != nil {
		return config, err // Return empty config if file doesn't exist
	}
	err = json.Unmarshal(data, &config)
	return config, err
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
