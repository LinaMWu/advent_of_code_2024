package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
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

	var (
		allLines string
		sumPart1 int64
		sumPart2 int64
	)

	// Combine all lines from input file into one string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		allLines += scanner.Text()
	}

	// Part 1
	sumPart1, err = calculateSum(allLines)
	if err != nil {
		logger.Error("failed to calculate sum for part 1", "error", err)
		os.Exit(1)
	}

	// Part 2
	// Split entire file by "do()" into sub strings
	subStrings := strings.Split(allLines, "do()")
	for _, subString := range subStrings {
		// Split each sub string by "don't()" into further sub strings
		subSubStrings := strings.Split(subString, "don't()")
		// Only need to calculate sum for the first portion of sub string that does not follow "don't()"
		subSum, err := calculateSum(subSubStrings[0])
		if err != nil {
			logger.Error("failed to calculate sub sum for part 2", "error", err, "sub sub string", subSubStrings[0])
			os.Exit(1)
		}

		sumPart2 += subSum
	}

	// Part 1 Solution
	logger.Info("result #1 is ready!", "sum part 1", sumPart1)

	// Part 2 Solution
	logger.Info("result #2 is ready!", "sum part 2", sumPart2)
}

func calculateSum(line string) (int64, error) {
	var sum int64

	// Regex to find mul(X,Y)
	reMul := regexp.MustCompile(`(?:mul\()(\d+)(?:,)(\d+)(?:\))`)
	matches := reMul.FindAllStringSubmatch(line, -1)

	for _, match := range matches {
		num1, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, err
		}

		num2, err := strconv.Atoi(match[2])
		if err != nil {
			return 0, err
		}

		sum += int64(num1) * int64(num2)
	}

	return sum, nil
}
