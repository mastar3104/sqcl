package readline_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestShiftTabPatchExists verifies that the Shift+Tab patch is applied
// to vendor/github.com/chzyer/readline/utils.go.
//
// If this test fails after running `go mod vendor`, you need to re-apply
// the patch by adding the following case to escapeExKey function:
//
//	case 'Z':
//		r = CharBackward // Shift+Tab
func TestShiftTabPatchExists(t *testing.T) {
	// Get the project root directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current file path")
	}
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")

	// Path to the readline utils.go
	utilsPath := filepath.Join(projectRoot, "vendor", "github.com", "chzyer", "readline", "utils.go")

	// Read the file
	content, err := os.ReadFile(utilsPath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", utilsPath, err)
	}

	// Check for the Shift+Tab patch
	if !strings.Contains(string(content), "case 'Z':") {
		t.Errorf("Shift+Tab patch not found in %s", utilsPath)
		t.Error("The patch maps Esc[Z (Shift+Tab) to CharBackward for reverse completion.")
		t.Error("Please re-apply the patch to vendor/github.com/chzyer/readline/utils.go")
	}
}
