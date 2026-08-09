package version

import (
	"testing"
)

func TestVersionIsSet(t *testing.T) {
	if Version == "" {
		t.Error("Version constant is empty")
	}
	if Version != "1.3.0" {
		t.Errorf("Version = %q, want %q", Version, "1.3.0")
	}
}
