package main

import (
	"bufio"
	"flag"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	increasing = "increasing"
	decreasing = "decreasing"
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
		part1valid = 0
		part2valid = 0
	)

	// Read input file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Fields(line)

		tmpLogger := logger.With(
			slog.String("line read", line),
			slog.Int("numbers per line", len(nums)),
			slog.Any("numbers found", nums),
		)

		// Check if report is safe
		safe, err := SafetyCheck(nums)
		if err != nil {
			tmpLogger.Error("unable to determine if report is safe", "error", err)
			os.Exit(1)
		}

		if safe {
			part1valid++
			part2valid++
		} else {
			for index := 0; index < len(nums); index++ {
				copyOfNums := make([]string, len(nums))
				copy(copyOfNums, nums)

				copyOfNums = append(copyOfNums[:index], copyOfNums[index+1:]...)

				safe, err := SafetyCheck(copyOfNums)
				if err != nil {
					tmpLogger.Error("unable to determine if report is safe", "error", err, "new numbers", copyOfNums)
					os.Exit(1)
				}

				if safe {
					part2valid++
					break
				}

			}
		}
	}

	// Part 1 Solution
	logger.Info("result #1 is ready!", "valid count", part1valid)

	// Part 2 Solution
	logger.Info("result #2 is ready!", "valid count", part2valid)
}

func SafetyCheck(nums []string) (bool, error) {
	previous := 0
	pattern := ""
	safe := false

	for index, num := range nums {

		level, err := strconv.Atoi(num)
		if err != nil {
			return false, err
		}

		// If this is the first level, save it and conitnue to the next level
		if index == 0 {
			previous = level
			continue
		}

		// If this is the second level, determine if the pattern is increasing or decreasing
		if index == 1 {
			if previous < level {
				pattern = increasing
			} else if previous > level {
				pattern = decreasing
			} else {
				// Invalid because it is not increasing or decreasing
				break
			}
		}

		// If this is the third level and on, check if it matches the established pattern
		if index > 1 {
			if pattern == increasing && previous > level {
				// Invalid because the pattern is increasing but the current level is decreasing
				break
			}

			if pattern == decreasing && previous < level {
				// Invalid because the pattern is decreasing but the current level is increasing
				break
			}
		}

		diff := math.Abs(float64(previous - level))
		if diff < 1 || diff > 3 {
			// Invalid because the difference is not between 1 and 3
			break
		}

		// If we are at the end, then this report passed all the checks
		if index == len(nums)-1 {
			safe = true
			break
		}

		previous = level
	}

	return safe, nil
}
