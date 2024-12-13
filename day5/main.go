package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"slices"
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

	pageOrderingRules := make(map[int][]int)
	var pageNumLists [][]int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		tmpLogger := logger.With(
			slog.String("line read", line),
		)

		// If line can be split by |, then it is a page ordering rule.
		subStrings := strings.Split(line, "|")
		if len(subStrings) == 2 {
			num1, err := strconv.Atoi(subStrings[0])
			if err != nil {
				tmpLogger.Error("failed to convert num1 from string to int", "error", err, "num1 string", subStrings[0])
				os.Exit(1)
			}

			num2, err := strconv.Atoi(subStrings[1])
			if err != nil {
				tmpLogger.Error("failed to convert num2 from string to int", "error", err, "num2 string", subStrings[1])
				os.Exit(1)
			}

			// num2 requires num1 to be listed before if both are listed
			pageOrderingRules[num2] = append(pageOrderingRules[num2], num1)
		}

		// If line can be split by , then it is a list of page numbers.
		nums := strings.Split(line, ",")
		if len(nums) > 1 {
			pageNumList := make([]int, len(nums))
			for index, num := range nums {
				pageNum, err := strconv.Atoi(num)
				if err != nil {
					tmpLogger.Error("failed to convert page number from string to int", "error", err, "page number string", num)
					os.Exit(1)
				}

				pageNumList[index] = pageNum
			}

			pageNumLists = append(pageNumLists, pageNumList)
		}
	}

	// Part 1
	var middlePageSum int
	// Part 2
	var fixedMiddlePageSum int

	// For each list of page numbers
	for _, pageNumList := range pageNumLists {
		correct, index, preReq := checkPageNumList(pageNumList, pageOrderingRules)

		if correct {
			middlePageSum += pageNumList[len(pageNumList)/2]
		} else {
			for {
				// To fix this list, first remove the pre-requisite page number from the list
				pageNumList = removeFromList(pageNumList, preReq)
				// Then add it right before the page number that requires it
				pageNumList = insertToList(pageNumList, preReq, index)

				correct, index, preReq = checkPageNumList(pageNumList, pageOrderingRules)
				if correct {
					fixedMiddlePageSum += pageNumList[len(pageNumList)/2]
					break
				}

			}
		}
	}

	// Part 1 Solution
	logger.Info("result #1 is ready!", "middle page sum", middlePageSum)

	// Part 2 Solution
	logger.Info("result #2 is ready!", "fixed middle page sum", fixedMiddlePageSum)
}

// checkPageNumList return if the list of page numbers is correct based on the ordering rules.
// If not, it also returns the index of the first incorrect page number and the missing pre-requisite page number.
func checkPageNumList(pageNumList []int, pageOrderingRules map[int][]int) (bool, int, int) {
	// For each page number in the list
	for index, pageNum := range pageNumList {
		// If the page number has ordering rulesm
		if rules, ok := pageOrderingRules[pageNum]; ok {
			// For each pre-requisite page number
			for _, preReqPage := range rules {
				// If the pre-requisite page number is in the list but no before the current page number
				if slices.Contains(pageNumList, preReqPage) && !slices.Contains(pageNumList[:index], preReqPage) {
					return false, index, preReqPage
				}
			}
		}
	}

	return true, 0, 0
}

// removeFromList removes an element from a slice.
func removeFromList[T comparable](l []T, item T) []T {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

// insertToList inserts an element into a slice at the specified index.
func insertToList[T any](array []T, value T, index int) []T {
	return append(array[:index], append([]T{value}, array[index:]...)...)
}
