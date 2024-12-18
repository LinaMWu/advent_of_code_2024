package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"strconv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Get input file name
	inputFlag := flag.String("input", "input/test_input.txt", "A file containing puzzle inputs.")
	flag.Parse()

	inputFileName := *inputFlag
	logger = logger.With(
		slog.String("inputFileName", inputFileName),
	)

	// Open input file
	file, err := os.Open(inputFileName)
	if err != nil {
		logger.Error("failed to open file", "error", err)
		os.Exit(1)
	}
	defer file.Close()

	var input string

	// Read input file into a single string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = scanner.Text()
	}

	// Generate file system based on input
	var fileSystem []string
	for i, char := range input {
		blockCount, err := strconv.Atoi(string(char))
		if err != nil {
			logger.Error("failed to convert block count from char to integer", "error", err, "char", string(char))
			os.Exit(1)
		}

		if i%2 == 0 {
			fileSystem = append(fileSystem, generateContent(strconv.Itoa(i/2), blockCount)...)
		} else {
			fileSystem = append(fileSystem, generateContent(".", blockCount)...)
		}
	}

	// Part 1
	squishedFileSystem := copySlice(fileSystem)
	for i := range squishedFileSystem {
		// If current character is ".", find the last number and swap them
		if squishedFileSystem[i] == "." {
			for j := len(squishedFileSystem) - 1; j > i; j-- {
				if squishedFileSystem[j] != "." {
					squishedFileSystem[i] = squishedFileSystem[j]
					squishedFileSystem[j] = "."
					break
				}
			}
		}
	}

	squishedCheckSum, err := checkSum(squishedFileSystem)
	if err != nil {
		logger.Error("failed to calculate check sum for squished file system", "error", err)
		os.Exit(1)
	}

	// Part 1 Solution
	logger.Info("result #1 is ready!", "squished check sum", squishedCheckSum)

	// Part 2
	reorgFileSystem := copySlice(fileSystem)
	index := len(reorgFileSystem)
	currFileID := ""
	blockCount := 0

	for index > 0 {
		index--

		// If the complete block of a fild ID has been identify, find the first opening and move the block
		if currFileID != "" && blockCount > 0 && reorgFileSystem[index] != currFileID {
			openingIndex := findOpening(reorgFileSystem, blockCount, index+1)
			if openingIndex != -1 {
				for i := 0; i < blockCount; i++ {
					reorgFileSystem[openingIndex+i] = currFileID
					reorgFileSystem[index+1+i] = "."
				}
			}

			currFileID = ""
			blockCount = 0
		}

		// If the current space is ".", skip it
		if reorgFileSystem[index] == "." {
			continue
		}

		// If the current space is a file ID, start counting the block
		if currFileID == "" {
			currFileID = reorgFileSystem[index]
			blockCount++
			continue
		}

		// If the current space is the same file ID, continue counting the block
		if reorgFileSystem[index] == currFileID {
			blockCount++
			continue
		}
	}

	reorgCheckSum, err := checkSum(reorgFileSystem)
	if err != nil {
		logger.Error("failed to calculate check sum for reorg file system", "error", err)
		os.Exit(1)
	}

	// Part 2 Solution
	logger.Info("result #2 is ready!", "reorg check sum", reorgCheckSum)
}

// generateContent generates a slice of strings with the given content at the given length
func generateContent(content string, count int) []string {
	result := make([]string, count)

	for i := range result {
		result[i] = content
	}

	return result
}

func copySlice[T any](s []T) []T {
	result := make([]T, len(s))
	copy(result, s)
	return result
}

func checkSum(fileSystem []string) (int, error) {
	var result int

	for i := range fileSystem {
		if fileSystem[i] != "." {
			id, err := strconv.Atoi(fileSystem[i])
			if err != nil {
				return 0, err
			}

			result += i * id
		}
	}

	return result, nil
}

// findOpening finds the first opening in the file system that fits the
// given block count, provided it is before the current index
func findOpening(fileSystem []string, blockCount, currentIndex int) int {
	for i := 0; i < len(fileSystem); i++ {
		if i >= currentIndex {
			break
		}

		if i+blockCount >= len(fileSystem) {
			break
		}

		if fileSystem[i] == "." {
			found := true

			for j := 0; j < blockCount; j++ {
				if fileSystem[i+j] != "." {
					found = false
					break
				}
			}

			if found {
				return i
			}
		}
	}

	return -1
}
