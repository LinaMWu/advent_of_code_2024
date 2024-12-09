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
		sumPart1 int64
		sumPart2 int64
	)

	// Read input file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		tmpLogger := logger.With(
			slog.String("line read", line),
		)

		subSum, err := calculateSum(line)
		if err != nil {
			tmpLogger.Error("failed to calculate sub sum", "error", err)
			os.Exit(1)
		}

		sumPart1 += subSum

		subStrings := strings.Split(line, "do()")
		for _, subString := range subStrings {
			subSubStrings := strings.Split(subString, "don't()")
			subsubSum, err := calculateSum(subSubStrings[0])
			if err != nil {
				tmpLogger.Error("failed to calculate sub sub sum", "error", err, "sub sub string", subSubStrings[0])
				os.Exit(1)
			}

			sumPart2 += subsubSum
		}
	}

	// Part 1 Solution
	logger.Info("result #1 is ready!", "sum part 1", sumPart1)

	// Part 2 Solution
	logger.Info("result #2 is ready!", "sum part 2", sumPart2) // incorrect
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
