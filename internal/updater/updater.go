package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/velo4705/polyglot/internal/ui"
)

const (
	githubAPIURL    = "https://api.github.com/repos/velo4705/polyglot/releases/latest"
	updateCheckFile = ".polyglot/last_update_check"
	checkInterval   = 24 * time.Hour
)

// Release represents a GitHub release
type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// Updater handles version checking and updates
type Updater struct {
	currentVersion string
	quiet          bool
}

// New creates a new updater
func New(currentVersion string, quiet bool) *Updater {
	return &Updater{
		currentVersion: currentVersion,
		quiet:          quiet,
	}
}

// CheckForUpdates checks if a new version is available
func (u *Updater) CheckForUpdates() (*Release, bool, error) {
	if !u.quiet {
		ui.Info("Checking for updates...")
	}

	// Fetch latest release info
	resp, err := http.Get(githubAPIURL)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("failed to fetch release info: status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, false, fmt.Errorf("failed to parse release info: %w", err)
	}

	// Compare versions
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := strings.TrimPrefix(u.currentVersion, "v")

	if latestVersion == currentVersion {
		if !u.quiet {
			ui.Success("You're running the latest version (%s)", currentVersion)
		}
		return &release, false, nil
	}

	if !u.quiet {
		ui.Info("New version available: %s (current: %s)", latestVersion, currentVersion)
	}

	return &release, true, nil
}

// Update downloads and installs the latest version
func (u *Updater) Update(release *Release) error {
	if !u.quiet {
		ui.Header("Updating Polyglot")
		fmt.Println()
	}

	// Determine binary name for current platform
	binaryName := u.getBinaryName()
	if binaryName == "" {
		return fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// Find matching asset
	var downloadURL string
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, binaryName) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	if !u.quiet {
		ui.Step("Downloading %s...", binaryName)
	}

	// Download new binary
	tmpFile, err := u.downloadBinary(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer os.Remove(tmpFile)

	if !u.quiet {
		ui.Success("Downloaded successfully")
	}

	// Get current binary path
	currentBinary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current binary path: %w", err)
	}
	currentBinary, err = filepath.EvalSymlinks(currentBinary)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	// Backup current binary
	backupPath := currentBinary + ".backup"
	if !u.quiet {
		ui.Step("Backing up current version...")
	}
	if err := u.copyFile(currentBinary, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Replace binary
	if !u.quiet {
		ui.Step("Installing new version...")
	}
	if err := u.replaceBinary(tmpFile, currentBinary); err != nil {
		// Restore backup on failure
		_ = u.copyFile(backupPath, currentBinary)
		os.Remove(backupPath)
		return fmt.Errorf("failed to install new version: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	if !u.quiet {
		ui.Success("Update complete!")
		fmt.Println()
		ui.Info("Polyglot has been updated to version %s", strings.TrimPrefix(release.TagName, "v"))
		ui.Dim("Run 'polyglot version' to verify")
	}

	return nil
}

// ShouldCheckForUpdates checks if it's time to check for updates
func (u *Updater) ShouldCheckForUpdates() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	checkFile := filepath.Join(home, updateCheckFile)
	info, err := os.Stat(checkFile)
	if err != nil {
		// File doesn't exist, should check
		return true
	}

	// Check if enough time has passed
	return time.Since(info.ModTime()) > checkInterval
}

// UpdateLastCheckTime updates the last check timestamp
func (u *Updater) UpdateLastCheckTime() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	checkFile := filepath.Join(home, updateCheckFile)
	dir := filepath.Dir(checkFile)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Touch file
	return os.WriteFile(checkFile, []byte(time.Now().Format(time.RFC3339)), 0644)
}

// getBinaryName returns the binary name for the current platform
func (u *Updater) getBinaryName() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	switch goos {
	case "linux":
		switch goarch {
		case "amd64":
			return "polyglot-linux-amd64"
		case "arm64":
			return "polyglot-linux-arm64"
		}
	case "darwin":
		switch goarch {
		case "amd64":
			return "polyglot-darwin-amd64"
		case "arm64":
			return "polyglot-darwin-arm64"
		}
	case "windows":
		if goarch == "amd64" {
			return "polyglot-windows-amd64.exe"
		}
	}

	return ""
}

// downloadBinary downloads a binary to a temporary file
func (u *Updater) downloadBinary(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "polyglot-update-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Download with progress
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	// Make executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

// copyFile copies a file
func (u *Updater) copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// replaceBinary replaces the current binary with a new one
func (u *Updater) replaceBinary(newBinary, currentBinary string) error {
	// On Windows, we can't replace a running executable
	// On Unix, we can replace it directly
	if runtime.GOOS == "windows" {
		// Move current binary to .old
		oldBinary := currentBinary + ".old"
		if err := os.Rename(currentBinary, oldBinary); err != nil {
			return err
		}

		// Copy new binary
		if err := u.copyFile(newBinary, currentBinary); err != nil {
			// Restore old binary
			_ = os.Rename(oldBinary, currentBinary)
			return err
		}

		// Schedule deletion of old binary
		go func() {
			time.Sleep(2 * time.Second)
			os.Remove(oldBinary)
		}()
	} else {
		// Unix: can replace directly
		if err := u.copyFile(newBinary, currentBinary); err != nil {
			return err
		}
	}

	return nil
}

// CheckForUpdatesInBackground checks for updates silently
func (u *Updater) CheckForUpdatesInBackground() {
	if !u.ShouldCheckForUpdates() {
		return
	}

	go func() {
		release, hasUpdate, err := u.CheckForUpdates()
		if err != nil {
			return
		}

		_ = u.UpdateLastCheckTime()

		if hasUpdate {
			fmt.Println()
			ui.Info("A new version of Polyglot is available: %s", strings.TrimPrefix(release.TagName, "v"))
			ui.Dim("Run 'polyglot update' to upgrade")
			fmt.Println()
		}
	}()
}

// UpdateViaPackageManager attempts to update via package manager
func (u *Updater) UpdateViaPackageManager() error {
	// Detect how Polyglot was installed
	if u.isInstalledViaHomebrew() {
		if !u.quiet {
			ui.Info("Detected Homebrew installation")
			ui.Step("Running: brew upgrade polyglot")
		}
		cmd := exec.Command("brew", "upgrade", "polyglot")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	if u.isInstalledViaApt() {
		if !u.quiet {
			ui.Info("Detected APT installation")
			ui.Step("Running: sudo apt update && sudo apt upgrade polyglot")
		}
		cmd := exec.Command("sudo", "apt", "update")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		cmd = exec.Command("sudo", "apt", "upgrade", "-y", "polyglot")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// No package manager detected — caller should use direct binary update
	return fmt.Errorf("no package manager detected")
}

// isInstalledViaHomebrew checks if installed via Homebrew
func (u *Updater) isInstalledViaHomebrew() bool {
	binary, err := os.Executable()
	if err != nil {
		return false
	}
	return strings.Contains(binary, "/Cellar/") || strings.Contains(binary, "/opt/homebrew/")
}

// isInstalledViaApt checks if installed via APT
func (u *Updater) isInstalledViaApt() bool {
	_, err := exec.LookPath("dpkg")
	if err != nil {
		return false
	}
	cmd := exec.Command("dpkg", "-l", "polyglot")
	return cmd.Run() == nil
}
