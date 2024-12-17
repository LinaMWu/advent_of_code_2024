package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
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
		part1Sum int
		part2Sum int
	)

	// Read input file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		tmpLogger := logger.With(
			slog.String("line read", line),
		)

		subStrings := strings.Split(line, ":")
		if len(subStrings) != 2 {
			tmpLogger.Error("invalid line read from file")
			os.Exit(1)
		}

		// Parse out target value
		target, err := strconv.Atoi(subStrings[0])
		if err != nil {
			tmpLogger.Error("failed to convert target value from string to int", "error", err, "target value string", subStrings[0])
			os.Exit(1)
		}

		// Parse out test values
		numStrings := strings.Fields(subStrings[1])
		nums := make([]int, len(numStrings))
		for i, numString := range numStrings {
			num, err := strconv.Atoi(numString)
			if err != nil {
				tmpLogger.Error("failed to convert num from string to int", "error", err, "num string", numString)
				os.Exit(1)
			}

			nums[i] = num
		}

		// Part 1
		if checkTargetPart1(target, nums, len(nums)-1) {
			part1Sum += target
		}

		// Part 2
		if checkTargetPart2(target, nums, len(nums)-1) {
			part2Sum += target
		}
	}

	// Part 1 Solution
	logger.Info("result #1 is ready!", "part 1 sum", part1Sum)

	// Part 2 Solution
	logger.Info("result #2 is ready!", "part 2 sum", part2Sum)
}

func checkTargetPart1(target int, nums []int, index int) bool {
	result := false

	// If we are at the last number
	if index == 0 {
		// Check if the target value is equal to the number
		if target == nums[index] {
			result = true
		}
	} else { // If we are not at the last number
		// Check if target value is larger than current number
		if target-nums[index] > 0 {
			// Subtract target value by current number and check the next number recursively
			result = result || checkTargetPart1(target-nums[index], nums, index-1)
		}

		// Check if target value is divisible by current number
		if target%nums[index] == 0 {
			// Divide target value by current number and check the next number recursively
			result = result || checkTargetPart1(target/nums[index], nums, index-1)
		}
	}

	return result
}

func checkTargetPart2(target int, nums []int, index int) bool {
	result := false

	// If we are at the last number
	if index == 0 {
		// Check if the target value is equal to the number
		if target == nums[index] {
			result = true
		}
	} else { // If we are not at the last number
		// Check if target value is larger than current number
		if target-nums[index] > 0 {
			// Subtract target value by current number and check the next number recursively
			result = result || checkTargetPart2(target-nums[index], nums, index-1)
		}

		// Check if target value is divisible by current number
		if target%nums[index] == 0 {
			// Divide target value by current number and check the next number recursively
			result = result || checkTargetPart2(target/nums[index], nums, index-1)
		}

		// Check if current number matches the lower digits of target value
		if target%padding(nums[index]) == nums[index] {
			// Remove the lower digits of target value and check the next number recursively
			result = result || checkTargetPart2(target/padding(nums[index]), nums, index-1)
		}
	}

	return result
}

// padding returns the smallest power of 10 that is larger than n
func padding(n int) int {
	var p int = 1

	for p <= n {
		p *= 10
	}

	return p
}
