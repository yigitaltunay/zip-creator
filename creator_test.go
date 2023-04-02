package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestZipSplitter(t *testing.T) {
	// Create some test files
	testDir := "."
	fileNames := make([]string, 0)
	for i := 0; i < 5; i++ {
		filename := fmt.Sprintf("testfile%d.txt", i)
		fileNames = append(fileNames, filename)
		file, _ := os.Create(filename)
		file.WriteString("This is a test file.")
		file.Close()
	}

	// Run the zip splitter
	Run(fileNames, 2)

	// Check that zip files were created
	zipFiles, err := filepath.Glob("*.zip")
	if err != nil {
		t.Fatalf("Error finding zip files: %s", err)
	}
	if len(zipFiles) != 3 {
		t.Fatalf("Expected 3 zip files, but got %d", len(zipFiles))
	}

	// Check that files are in the zip files
	for _, zipFile := range zipFiles {
		reader, err := zip.OpenReader(zipFile)
		if err != nil {
			t.Fatalf("Error opening zip file %s: %s", zipFile, err)
		}
		defer reader.Close()
		if len(reader.File) != 2 {
			t.Fatalf("Expected 2 files per zip file, but got %d", len(reader.File))
		}
		for _, file := range reader.File {
			if file.FileInfo().IsDir() {
				continue
			}
			f, err := file.Open()
			if err != nil {
				t.Fatalf("Error opening file %s in zip file %s: %s", file.Name, zipFile, err)
			}
			defer f.Close()
			content, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatalf("Error reading content of file %s in zip file %s: %s", file.Name, zipFile, err)
			}
			expectedContent := "This is a test file."
			if string(content) != expectedContent {
				t.Fatalf("Expected file content to be '%s', but got '%s'", expectedContent, string(content))
			}
		}
	}

	// Cleanup test files and zip files
	os.RemoveAll(testDir)
	for _, zipFile := range zipFiles {
		os.Remove(zipFile)
	}
}
