package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// This is the major version
const MajorVersion = "0"

// This variable will be overwritten at compile time
var Version = "dev"

// GitHubCommit simplified struct for parsing the JSON response
type GitHubCommit struct {
	SHA string `json:"sha"`
}

// ****************************************************************************
// GetDisplayVersion()
// ****************************************************************************
func GetDisplayVersion() string {
	if Version == "dev" {
		return MajorVersion + ".x-dev"
	}
	return Version
}

// ****************************************************************************
// fetchRemoteHash()
// ****************************************************************************
func fetchRemoteHash() (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	// We use the GitHub API to get the latest commit from the main branch
	resp, err := client.Get("https://api.github.com/repos/jplozf/pingo/commits/main")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result GitHubCommit
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// GitHub returns a long 40-character hash; we take the first 7 for a short hash
	if len(result.SHA) >= 7 {
		return result.SHA[:7], nil
	}
	return result.SHA, nil
}

// ****************************************************************************
// checkForUpdates()
// ****************************************************************************
func checkForUpdates(w fyne.Window) {
	remoteHash, err := fetchRemoteHash()
	if err != nil {
		// Silently fail or log error to avoid bothering user on offline mode
		return
	}

	// We assume your 'Version' string ends with the hash (e.g., "0.5-abcdef")
	// We check if the remote hash is present in our local Version string
	if !strings.Contains(Version, remoteHash) {
		dialog.ShowInformation("Update Available",
			"A new version is available on GitHub!\nRemote Hash: "+remoteHash,
			w)
	}
}
