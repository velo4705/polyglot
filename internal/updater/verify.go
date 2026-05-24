package updater

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// VerifyChecksum downloads the .sha256 sidecar for the asset and verifies the
// SHA256 hash of the local file at binaryPath.
// checksumURL is the URL of the checksum file (e.g. the browser_download_url
// for the "<binaryName>.sha256" asset).
func VerifyChecksum(binaryPath, checksumURL string) error {
	// Download checksum file
	resp, err := http.Get(checksumURL)
	if err != nil {
		return fmt.Errorf("downloading checksum: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("downloading checksum: status %d", resp.StatusCode)
	}

	checksumData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading checksum: %w", err)
	}

	// Checksum files are typically "<hex>  <filename>" (sha256sum format)
	expectedHex := strings.Fields(string(checksumData))[0]

	// Compute SHA256 of downloaded binary
	f, err := os.Open(binaryPath)
	if err != nil {
		return fmt.Errorf("opening binary for checksum: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("hashing binary: %w", err)
	}

	actualHex := hex.EncodeToString(h.Sum(nil))

	if !strings.EqualFold(actualHex, expectedHex) {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedHex, actualHex)
	}

	return nil
}

// VerifyGPGSignature downloads the .asc detached signature for the asset and
// verifies it using the local gpg binary.
// signatureURL is the URL of the .asc file.
func VerifyGPGSignature(binaryPath, signatureURL string) error {
	// Check that gpg is available
	if _, err := exec.LookPath("gpg"); err != nil {
		return fmt.Errorf("gpg not found in PATH; cannot verify signature (install gpg to enable verification)")
	}

	// Download .asc signature to a temp file
	resp, err := http.Get(signatureURL)
	if err != nil {
		return fmt.Errorf("downloading signature: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("downloading signature: status %d", resp.StatusCode)
	}

	tmpSig, err := os.CreateTemp("", "polyglot-sig-*.asc")
	if err != nil {
		return fmt.Errorf("creating temp signature file: %w", err)
	}
	defer os.Remove(tmpSig.Name())
	defer tmpSig.Close()

	if _, copyErr := io.Copy(tmpSig, resp.Body); copyErr != nil {
		return fmt.Errorf("writing signature: %w", copyErr)
	}
	tmpSig.Close()

	// Run: gpg --verify <sig> <binary>
	cmd := exec.Command("gpg", "--verify", tmpSig.Name(), binaryPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("GPG verification failed: %s", strings.TrimSpace(string(out)))
	}

	return nil
}
