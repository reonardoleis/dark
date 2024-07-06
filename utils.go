package main

import "fmt"

func printIntGrid(grid [][]int) {
	for _, line := range grid {
		for _, value := range line {
			fmt.Print(value, " ")
		}
		fmt.Println("")
	}
}
