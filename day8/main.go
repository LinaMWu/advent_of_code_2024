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

	// Part 1
	part1Matrix := copyMatrix(matrix)
	for i, row := range matrix {
		for j, col := range row {
			if col != "." {
				findAntinodesPart1(matrix, part1Matrix, i, j, col)
			}
		}
	}

	// Part 1 Solution
	logger.Info("result #1 is ready!", "antinodes count", countAntinodes(part1Matrix))

	// Part 2
	part2Matrix := copyMatrix(matrix)
	for i, row := range matrix {
		for j, col := range row {
			if col != "." {
				findAntinodesPart2(matrix, part2Matrix, i, j, col)
			}
		}
	}

	// Part 2 Solution
	logger.Info("result #2 is ready!", "antinodes count", countAntinodes(part2Matrix))
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

func findAntinodesPart1(refMatrix, markingMatrix [][]string, currI, currJ int, freq string) {
	for i, row := range refMatrix {
		for j, col := range row {
			// Skip current node
			if i == currI && j == currJ {
				continue
			}

			// Check for resonant frequency
			if col == freq {
				// Calculate offset
				offsetI := currI - i
				offsetJ := currJ - j

				// Apply offset to location of the resonant frequency
				if i-offsetI >= 0 && i-offsetI < len(markingMatrix) && j-offsetJ >= 0 && j-offsetJ < len(markingMatrix[0]) {
					markingMatrix[i-offsetI][j-offsetJ] = "#"
				}
			}
		}
	}
}

func findAntinodesPart2(refMatrix, markingMatrix [][]string, currI, currJ int, freq string) {
	for i, row := range refMatrix {
		for j, col := range row {
			// Skip current node
			if i == currI && j == currJ {
				continue
			}

			// Check for resonant frequency
			if col == freq {
				// Calculate offset
				offsetI := currI - i
				offsetJ := currJ - j

				// Find location of the antinode
				antinodeI := i
				antinodeJ := j
				for {
					// Check if antinode is within bounds
					if antinodeI < 0 || antinodeI >= len(markingMatrix) || antinodeJ < 0 || antinodeJ >= len(markingMatrix[0]) {
						break
					}

					// Mark antinode
					markingMatrix[antinodeI][antinodeJ] = "#"

					// Check for the next antinode location
					antinodeI -= offsetI
					antinodeJ -= offsetJ
				}
			}
		}
	}
}

func countAntinodes(matrix [][]string) int {
	count := 0
	for _, row := range matrix {
		for _, col := range row {
			if col == "#" {
				count++
			}
		}
	}

	return count
}
