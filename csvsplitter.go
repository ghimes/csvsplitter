package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var inputFileName *string
var outputFileName *string
var linesPerFile *int

func readArguments() {

	inputFileName = flag.String("input-file", "data/example.csv", "a string")
	outputFileName = flag.String("output-file", *inputFileName, "a string")
	linesPerFile = flag.Int("lines", 100000, "an int")

	flag.Parse()

	if *linesPerFile <= 0 {
		log.Fatal("Number of lines per file should be higher than 0")
		os.Exit(1)
	}

	if !fileExists(*inputFileName) {
		log.Fatal("Input file does not exist!")
		os.Exit(1)
	}
}

func createDirectory(directoryPath string) {
	//choose your permissions well
	pathErr := os.MkdirAll(directoryPath, 0777)

	//check if you need to panic, fallback or report
	if pathErr != nil {
		fmt.Println(pathErr)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func writeCsv(name string, data [][]string) {

	filePart, err := os.Create(name)
	checkError("Cannot create file", err)
	defer func() {
		_ = filePart.Close()
	}()

	writer := csv.NewWriter(filePart)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}
}

func getFileNameAndPath(file string) (dir string, fileName string) {
	dir, fileName = filepath.Split(file)
	return dir, fileName
}

func main() {

	start := time.Now()

	readArguments()

	outputFileDir, outputFileName := getFileNameAndPath(*outputFileName)

	if _, err := os.Stat(outputFileDir); os.IsNotExist(err) {
		createDirectory(outputFileDir)
	}

	fmt.Printf("File to be splitted: %s\n", *inputFileName)
	fmt.Printf("Number of lines per file: %d\n", *linesPerFile)

	// Load a csv file.
	f, err := os.Open(*inputFileName)
	checkError("Cannot open the file", err)

	defer func() {
		_ = f.Close()
	}()

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))

	currentLineNumber := -1
	fileCounter := 0

	var dataChunk [][]string
	var header []string

	for {
		line, endOfFile := r.Read()

		if currentLineNumber == -1 {
			header = line
		}

		if (dataChunk == nil) && (header != nil) && (currentLineNumber > -1) {
			dataChunk = append(dataChunk, header)
		}

		// Stop at EOF.
		if endOfFile == io.EOF {
			if len(dataChunk) > 1 {
				// write to file
				filePartName := fmt.Sprintf("%s%d_%s", outputFileDir, fileCounter, outputFileName)
				writeCsv(filePartName, dataChunk)
			}
			break
		}

		dataChunk = append(dataChunk, line)
		currentLineNumber++

		if currentLineNumber >= *linesPerFile {
			// write to file
			filePartName := fmt.Sprintf("%s%d_%s", outputFileDir, fileCounter, outputFileName)
			writeCsv(filePartName, dataChunk)

			fileCounter++
			currentLineNumber = 0
			dataChunk = nil
		}
	}
	elapsed := time.Since(start)
	log.Printf("Script took %s", elapsed)
}
