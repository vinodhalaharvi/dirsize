package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0K"},
		{1023, "0K"},
		{1024, "1K"},
		{1025, "1K"},
		{10240, "10K"},
	}

	for _, test := range tests {
		result := formatBytes(test.input)
		if result != test.expected {
			t.Errorf("formatBytes(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestFormatHuman(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0B"},
		{1023, "1023B"},
		{1024, "1.0K"},
		{1048576, "1.0M"},
		{1073741824, "1.0G"},
		{1099511627776, "1.0T"},
	}

	for _, test := range tests {
		result := formatHuman(test.input)
		if result != test.expected {
			t.Errorf("formatHuman(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestCalculateSizeWithErrors(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dirsize-test-errors")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a directory with no read permissions
	noAccessDir := filepath.Join(tempDir, "no_access")
	if err := os.Mkdir(noAccessDir, 0); err != nil {
		t.Fatalf("Failed to create no_access dir: %v", err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	calculateSize(tempDir, true, formatBytes)

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Error accessing") {
		t.Errorf("Expected error message for accessing restricted directory, but got none")
	}
}

func TestCalculateSize(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dirsize-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	createTestFiles(t, tempDir)

	tests := []struct {
		name      string
		path      string
		recursive bool
		expected  int64
	}{
		{"NonRecursiveRoot", tempDir, false, 19},
		{"RecursiveRoot", tempDir, true, 19},
		{"SingleFile", filepath.Join(tempDir, "file1.txt"), false, 13},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			result := calculateSize(test.path, test.recursive, formatBytes)

			// Restore stdout
			w.Close()
			os.Stdout = old

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if result != test.expected {
				t.Errorf("calculateSize(%s, %v) = %d; want %d", test.path, test.recursive, result, test.expected)
			}

			if test.recursive {
				if !strings.Contains(output, "file2.txt") {
					t.Errorf("Expected output to contain 'file2.txt' in recursive mode, but it didn't")
				}
			} else {
				if strings.Contains(output, "file2.txt") {
					t.Errorf("Expected output not to contain 'file2.txt' in non-recursive mode, but it did")
				}
			}
		})
	}
}

func createTestFiles(t *testing.T, root string) {
	if err := os.Mkdir(filepath.Join(root, "subdir1"), 0755); err != nil {
		t.Fatalf("Failed to create subdir1: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "subdir2"), 0755); err != nil {
		t.Fatalf("Failed to create subdir2: %v", err)
	}

	if err := os.WriteFile(filepath.Join(root, "file1.txt"), []byte("Hello, World!"), 0644); err != nil {
		t.Fatalf("Failed to create file1.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "subdir1", "file2.txt"), []byte("OpenAI"), 0644); err != nil {
		t.Fatalf("Failed to create file2.txt: %v", err)
	}
}
