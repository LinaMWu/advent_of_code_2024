package main

import (
	"bufio"
	"flag"
	"log/slog"
	"math"
	"os"
	"sort"
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
		index int
		list1 []int
		list2 []int
	)

	// Read input file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		index++
		line := scanner.Text()
		nums := strings.Fields(line)

		tmpLogger := logger.With(
			slog.Int("line number", index),
			slog.String("line read", line),
			slog.Int("numbers per line", len(nums)),
			slog.Any("numbers found", nums),
		)

		if len(nums) != 2 {
			tmpLogger.Error("encountered invalid input")
			os.Exit(1)
		}

		int1, err := strconv.Atoi(nums[0])
		if err != nil {
			tmpLogger.Error("unable to convert number 1 from string to int", "error", err, "number1", nums[0])
			os.Exit(1)
		}

		int2, err := strconv.Atoi(nums[1])
		if err != nil {
			tmpLogger.Error("unable to convert number 2 from string to int", "error", err, "number2", nums[1])
			os.Exit(1)
		}

		list1 = append(list1, int1)
		list2 = append(list2, int2)
	}

	logger = logger.With(
		slog.Int("list1 lenghth", len(list1)),
		slog.Int("list2 lenghth", len(list2)),
	)

	// Confirm two lists are the same length
	if len(list1) != len(list2) {
		logger.Error("list1 and list2 are not the same length")
		os.Exit(1)
	}

	// Sort the lists in ascending order
	sort.Ints(list1)
	sort.Ints(list2)

	// Sum up the absolute difference between the two lists
	var sum int
	for i := 0; i < len(list1); i++ {
		sum += int(math.Abs(float64(list1[i] - list2[i])))
	}

	// Part 1 Solution
	logger.Info("result #1 is ready!", "sum", sum)

	var (
		score     int64
		index2    int
		prevScore int64
	)

	// Calculate similarity score between the two lists
	for i := 0; i < len(list1); i++ {
		count := 0

		// If the current number is the same as the previous number, add the previous score to the current score and move on
		if i != 0 && list1[i] == list1[i-1] {
			score += prevScore
			continue
		}

		// If the current number is larger than the current number in the second list, move on to the next number in the second list
		for list1[i] > list2[index2] && index2 < len(list2)-1 {
			index2++
		}

		// If the current number is the same as the current number in the second list, count the number of times the number appears in the second list
		for list1[i] == list2[index2] && index2 < len(list2)-1 {
			count++
			index2++
		}

		// Save the current score and add it to the total score
		prevScore = int64(list1[i] * count)
		score += prevScore

		// If we reach the end of the second list, break the loop
		if index2 >= len(list2) {
			break
		}
	}

	// Part 2 Solution
	logger.Info("result #2 is ready!", "score", score)
}
