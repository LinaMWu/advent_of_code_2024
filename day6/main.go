package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
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

	var matrix [][]string

	// Read input file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]string, len(line))
		for i, char := range line {
			row[i] = string(char)
		}

		matrix = append(matrix, row)
	}

	// Find the guard in the map
	guardFound, startingI, startingJ, startingGuard := findGuard(matrix)
	if !guardFound {
		logger.Error("guard not found in map")
		os.Exit(1)
	}

	// Part 1
	part1Matrix := copyMatrix(matrix)
	i, j, guard := startingI, startingJ, startingGuard
	for {
		i, j, guard = moveGuard(part1Matrix, i, j, guard)
		if i == -1 && j == -1 {
			break
		}
	}

	// Part 1 Solution
	logger.Info("result #1 is ready!", "space count", countSpaces(part1Matrix))

	// Part 2
	stuckCount := 0
	maxSteps := len(matrix) * len(matrix[0])

	for row := 0; row < len(matrix); row++ {
		for col := 0; col < len(matrix[row]); col++ {
			part2Matrix := copyMatrix(matrix)

			// Turn the current cell into a wall if possible
			if matrix[row][col] != "." {
				continue
			} else {
				part2Matrix[row][col] = "#"
			}

			i, j, guard := startingI, startingJ, startingGuard
			steps := 0

			for {
				i, j, guard = moveGuard(part2Matrix, i, j, guard)
				if i == -1 && j == -1 {
					break
				}

				steps++

				// If the guard has taken more steps than the maximum possible steps, it's stuck
				if steps > maxSteps {
					stuckCount++
					break
				}
			}
		}
	}

	// Part 2 Solution
	logger.Info("result #2 is ready!", "stuck count", stuckCount)

}

func findGuard(matrix [][]string) (bool, int, int, string) {
	for i, row := range matrix {
		for j, cell := range row {
			if cell == "^" || cell == "v" || cell == "<" || cell == ">" {
				return true, i, j, cell
			}
		}
	}

	return false, -1, -1, ""
}

func moveGuard(matrix [][]string, i, j int, guard string) (int, int, string) {
	switch guard {
	case "^":
		// Move guard up until it hits a wall
		for index := i; index >= 0; index-- {
			if matrix[index][j] == "#" {
				return index + 1, j, ">"
			} else {
				matrix[index][j] = "X"
			}
		}
	case "v":
		// Move guard down until it hits a wall
		for index := i; index < len(matrix); index++ {
			if matrix[index][j] == "#" {
				return index - 1, j, "<"
			} else {
				matrix[index][j] = "X"
			}
		}
	case "<":
		// Move guard left until it hits a wall
		for index := j; index >= 0; index-- {
			if matrix[i][index] == "#" {
				return i, index + 1, "^"
			} else {
				matrix[i][index] = "X"
			}
		}
	case ">":
		// Move guard right until it hits a wall
		for index := j; index < len(matrix[i]); index++ {
			if matrix[i][index] == "#" {
				return i, index - 1, "v"
			} else {
				matrix[i][index] = "X"
			}
		}
	}

	return -1, -1, ""
}

func countSpaces(matrix [][]string) int {
	count := 0
	for _, row := range matrix {
		for _, cell := range row {
			if cell == "X" {
				count++
			}
		}
	}

	return count
}

func copyMatrix[T any](m [][]T) [][]T {
	copyMatrix := make([][]T, len(m))
	for index, row := range m {
		copyRow := make([]T, len(row))
		copy(copyRow, row)
		copyMatrix[index] = copyRow
	}
	return copyMatrix
}
