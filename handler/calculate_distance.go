package handler

import (
	"fmt"
)

type Tree struct {
	X, Y, Height int
}

// This function calculates the total distance a drone travels over a rectangular estate with trees. ...
func CalculateDroneDistance(length, width int, trees []Tree) int {
	// Fast tree height lookup
	treeMap := make(map[[2]int]int)
	for _, t := range trees {
		treeMap[[2]int{t.X, t.Y}] = t.Height
	}

	totalDistance := 0
	currentHeight := 0 // Start at ground level

	for y := 1; y <= width; y++ {
		var xStart, xEnd, step int
		if y%2 == 1 {
			// East
			xStart, xEnd, step = 1, length+1, 1
		} else {
			// West
			xStart, xEnd, step = length, 0, -1
		}

		for x := xStart; x != xEnd; x += step {
			targetHeight := 1 // Default for empty plot
			if h, ok := treeMap[[2]int{x, y}]; ok {
				targetHeight = h + 1
			}

			// Adjust vertical
			verticalChange := abs(currentHeight - targetHeight)
			totalDistance += verticalChange
			currentHeight = targetHeight

			// Only add horizontal movement if not at the last plot overall
			isLastPlot := (x == xEnd-step && y == width)
			if !isLastPlot {
				// Move to next plot
				totalDistance += 10
			}

			fmt.Printf("Moving to (%d, %d) with height %d, total distance: %d\n", x, y, targetHeight, totalDistance)
		}

		// Between rows: move north 10m (no vertical movement)
		if y < width {
			totalDistance += 10
		}
	}

	// Final descent to ground
	totalDistance += currentHeight

	return totalDistance
}

// This function calculates the total distance a drone travels over a rectangular estate with trees. ...
func MaxDistanceDrone(length, width int, trees []Tree, maxDistance int) []int {
	// Fast tree height lookup
	treeMap := make(map[[2]int]int)
	for _, t := range trees {
		treeMap[[2]int{t.X, t.Y}] = t.Height
	}

	totalDistance := 0
	currentHeight := 0 // Start at ground level

	lastCoordinates := []int{0, 0}

	for y := 1; y <= width; y++ {
		var xStart, xEnd, step int
		if y%2 == 1 {
			// East
			xStart, xEnd, step = 1, length+1, 1
		} else {
			// West
			xStart, xEnd, step = length, 0, -1
		}

		for x := xStart; x != xEnd; x += step {
			targetHeight := 1 // Default for empty plot
			if h, ok := treeMap[[2]int{x, y}]; ok {
				targetHeight = h + 1
			}

			// Adjust vertical
			verticalChange := abs(currentHeight - targetHeight)
			totalDistance += verticalChange
			currentHeight = targetHeight

			// Only add horizontal movement if not at the last plot overall
			isLastPlot := (x == xEnd-step && y == width)
			if !isLastPlot {
				// Move to next plot
				totalDistance += 10
			}
			if totalDistance > maxDistance {
				fmt.Printf("Max distance reached at (%d, %d) with height %d, total distance: %d\n", x, y, targetHeight, totalDistance)
				return []int{x, y}
			}

			fmt.Printf("Moving to (%d, %d) with height %d, total distance: %d\n", x, y, targetHeight, totalDistance)
			lastCoordinates = []int{x, y}
		}
		// Between rows: move north 10m (no vertical movement)
		if y < width {
			totalDistance += 10
		}

	}
	return lastCoordinates
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
