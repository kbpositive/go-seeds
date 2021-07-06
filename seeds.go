package main

import (
	"fmt"
	"strconv"
	"strings"
)

// define moves
var moves = [8][2]int{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
	{1, 0},
	{1, -1},
	{0, -1}}

// update grid function
func update(grid map[string]int) map[string]int {
	// copy grid
	var decoy = make(map[string]int)
	for i, v := range grid {
		decoy[i] = v
	}

	var memo = make(map[string]int)
	for key, _ := range grid {
		var stack = []string{key}
		for stack != nil {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if _, found := memo[cur]; found != true {
				memo[cur] = 1

				// valid moves
				var valid []string
				for _, move := range moves {
					row := strings.Split(cur, ",")[0]
					col := strings.Split(cur, ",")[1]
					introw, _ := strconv.Atoi(row)
					intcol, _ := strconv.Atoi(col)
					introw += move[0]
					intcol += move[1]

					valid = append(valid, strconv.Itoa(introw)+","+strconv.Itoa(intcol))
				}

				// valid children
				var children []int
				for _, move := range valid {
					if _, found := grid[move]; found {
						children = append(children, grid[move])
					}
				}

				// add children if cell is alive
				if _, found := grid[cur]; found {
					for _, child := range valid {
						stack = append(stack, child)
					}
				}

				// game rules B2/S
				if _, found := grid[cur]; found != true && len(children) == 2 {
					decoy[cur] = 1
				} else if _, found := decoy[cur]; found {
					delete(decoy, cur)
				}

			}

		}
	}
	return decoy
}

// render grid function
// block cellular automaton

func main() {

	fmt.Println(moves)
}
