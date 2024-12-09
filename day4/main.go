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
	var XMAScount int
	for row := 0; row < len(matrix); row++ {
		for col := 0; col < len(matrix[row]); col++ {
			if matrix[row][col] == "X" {
				XMAScount += checkXMAS(row, col, matrix)
			}
		}
	}

	// Part 1 Solution
	logger.Info("result #1 is ready!", "XMAS count", XMAScount)

	// Part 2
	var MAScount int
	for row := 0; row < len(matrix); row++ {
		for col := 0; col < len(matrix[row]); col++ {
			if matrix[row][col] == "A" && checkMAS(row, col, matrix) {
				MAScount++
			}
		}
	}

	// Part 2 Solution
	logger.Info("result #2 is ready!", "X-MAS count", MAScount)
}

func checkXMAS(row int, col int, matrix [][]string) int {
	var count int

	if checkHorizontalForward(row, col, matrix) {
		count++
	}

	if checkHorizontalBackward(row, col, matrix) {
		count++
	}

	if checkVerticalUpward(row, col, matrix) {
		count++
	}

	if checkVerticalDownward(row, col, matrix) {
		count++
	}

	if checkDiagonalUpwardRight(row, col, matrix) {
		count++
	}

	if checkDiagonalDownwardRight(row, col, matrix) {
		count++
	}

	if checkDiagonalUpwardLeft(row, col, matrix) {
		count++
	}

	if checkDiagonalDownwardLeft(row, col, matrix) {
		count++
	}

	return count
}

func checkHorizontalForward(row int, col int, matrix [][]string) bool {
	if col+3 >= len(matrix[row]) {
		return false
	}

	if matrix[row][col+1] == "M" && matrix[row][col+2] == "A" && matrix[row][col+3] == "S" {
		return true
	}

	return false
}

func checkHorizontalBackward(row int, col int, matrix [][]string) bool {
	if col-3 < 0 {
		return false
	}

	if matrix[row][col-1] == "M" && matrix[row][col-2] == "A" && matrix[row][col-3] == "S" {
		return true
	}

	return false
}

func checkVerticalUpward(row int, col int, matrix [][]string) bool {
	if row-3 < 0 {
		return false
	}

	if matrix[row-1][col] == "M" && matrix[row-2][col] == "A" && matrix[row-3][col] == "S" {
		return true
	}

	return false
}

func checkVerticalDownward(row int, col int, matrix [][]string) bool {
	if row+3 >= len(matrix) {
		return false
	}

	if matrix[row+1][col] == "M" && matrix[row+2][col] == "A" && matrix[row+3][col] == "S" {
		return true
	}

	return false
}

func checkDiagonalUpwardRight(row int, col int, matrix [][]string) bool {
	if row-3 < 0 || col+3 >= len(matrix[row]) {
		return false
	}

	if matrix[row-1][col+1] == "M" && matrix[row-2][col+2] == "A" && matrix[row-3][col+3] == "S" {
		return true
	}

	return false
}

func checkDiagonalDownwardRight(row int, col int, matrix [][]string) bool {
	if row+3 >= len(matrix) || col+3 >= len(matrix[row]) {
		return false
	}

	if matrix[row+1][col+1] == "M" && matrix[row+2][col+2] == "A" && matrix[row+3][col+3] == "S" {
		return true
	}

	return false
}

func checkDiagonalUpwardLeft(row int, col int, matrix [][]string) bool {
	if row-3 < 0 || col-3 < 0 {
		return false
	}

	if matrix[row-1][col-1] == "M" && matrix[row-2][col-2] == "A" && matrix[row-3][col-3] == "S" {
		return true
	}

	return false
}

func checkDiagonalDownwardLeft(row int, col int, matrix [][]string) bool {
	if row+3 >= len(matrix) || col-3 < 0 {
		return false
	}

	if matrix[row+1][col-1] == "M" && matrix[row+2][col-2] == "A" && matrix[row+3][col-3] == "S" {
		return true
	}

	return false
}

func checkMAS(row int, col int, matrix [][]string) bool {
	if row-1 < 0 || row+1 >= len(matrix) || col-1 < 0 || col+1 >= len(matrix[row]) {
		return false
	}

	if ((matrix[row-1][col-1] == "M" && matrix[row+1][col+1] == "S") || (matrix[row-1][col-1] == "S" && matrix[row+1][col+1] == "M")) && ((matrix[row-1][col+1] == "M" && matrix[row+1][col-1] == "S") || (matrix[row-1][col+1] == "S" && matrix[row+1][col-1] == "M")) {
		return true
	}

	return false
}
