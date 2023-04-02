package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	defaultSplitSize    = 0
	defaultFileLocation = "."
	defaultFileType     = "txt"
)

func main() {
	fileLocation := getInput("Write the file location ( put period to get current directory ) :", defaultFileLocation)
	fileType := getInput("Write the file type (example write 'txt') :", defaultFileType)
	files, err := ioutil.ReadDir(fileLocation)
	if err != nil {
		log.Fatal(err)
	}
	filesToZip := filterFiles(files, fileType)
	fmt.Printf("Total number of %s count: %d\n\n", fileType, len(filesToZip))
	splitSize := getSplitSize()
	Run(filesToZip, splitSize)
}

func getInput(prompt, defaultValue string) string {
	color.Green(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input. Please try again", err)
		return defaultValue
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func filterFiles(files []fs.FileInfo, fileType string) []string {
	var filteredFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "."+fileType) {
			filteredFiles = append(filteredFiles, filepath.Join(defaultFileLocation, file.Name()))
		}
	}
	return filteredFiles
}

func getSplitSize() int {
	input := getInput("How many parts will it be divided into? :", strconv.Itoa(defaultSplitSize))
	splitSize, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input. Using default value.")
		return defaultSplitSize
	}
	if splitSize <= 0 {
		fmt.Println("Invalid input. Using default value.")
		return defaultSplitSize
	}
	return splitSize
}

func Run(files []string, splitSize int) {
	chunks := chunkSlice(files, splitSize)
	for i, chunk := range chunks {
		zipFileName := fmt.Sprintf("%d.zip", i)
		createZipFile(zipFileName, chunk)
	}
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func createZipFile(zipFileName string, files []string) {
	fmt.Println("Creating zip archive:", zipFileName)
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	for _, file := range files {
		addFileToZip(file, zipWriter)
	}
	fmt.Println("Zip archive created:", zipFileName)
}

func addFileToZip(file string, zipWriter *zip.Writer) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	zipFile, err := zipWriter.Create(filepath.Base(file))
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(zipFile, f); err != nil {
		panic(err)
	}
}
