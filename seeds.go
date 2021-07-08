package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

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
func update(grid *map[string]int, grid_memo *map[string]int, quadrant int) {
	defer wg.Done()
	var memo = make(map[string]int)

	for key := range *grid {
		var stack = []string{key}
		for len(stack) > 0 {
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

				// live children
				var children []int
				for _, move := range valid {
					if _, found := (*grid)[move]; found {
						children = append(children, (*grid)[move])
					}
				}

				// add children if cell is alive
				if _, found := (*grid)[cur]; found {
					for _, child := range valid {
						y := strings.Split(child, ",")[0]
						x := strings.Split(child, ",")[1]
						yval, _ := strconv.Atoi(x)
						xval, _ := strconv.Atoi(y)
						if quadrant == 0 && yval <= 375 && xval <= 375 {
							stack = append(stack, child)
						} else if quadrant == 1 && yval <= 375 && xval >= 375 {
							stack = append(stack, child)
						} else if quadrant == 2 && yval >= 375 && xval <= 375 {
							stack = append(stack, child)
						} else if quadrant == 3 && yval >= 375 && xval >= 375 {
							stack = append(stack, child)
						}
					}
				}

				// game rules B2/S
				_, found := (*grid)[cur]
				if found != true && len(children) == 2 {
					(*grid_memo)[cur] = 1
				} else if found {
					delete(*grid_memo, cur)
				}

			}

		}
	}
}

// render grid
func render(grid *map[string]int, frames int, dim int) {
	var palette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0xff, 0xff},
		color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0xff, 0xff},
		color.RGBA{0xff, 0xff, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff},
	}
	var images []*image.Paletted
	var delays []int

	for step := 0; step < frames; step++ {
		img := image.NewPaletted(image.Rect(0, 0, dim, dim), palette)
		images = append(images, img)
		delays = append(delays, 0)

		for i, _ := range *grid {
			row := strings.Split(i, ",")[0]
			col := strings.Split(i, ",")[1]
			rowval, _ := strconv.Atoi(row)
			colval, _ := strconv.Atoi(col)
			img.Set(rowval, colval, color.RGBA{255, 255, 255, 255})
		}

		var grid_q1 = make(map[string]int)
		var grid_q2 = make(map[string]int)
		var grid_q3 = make(map[string]int)
		var grid_q4 = make(map[string]int)
		for i, v := range *grid {
			row := strings.Split(i, ",")[0]
			col := strings.Split(i, ",")[1]
			rowval, _ := strconv.Atoi(row)
			colval, _ := strconv.Atoi(col)
			if rowval <= 375 && colval <= 375 {
				grid_q1[i] = v
			} else if rowval <= 375 && colval >= 375 {
				grid_q2[i] = v
			} else if rowval >= 375 && colval <= 375 {
				grid_q3[i] = v
			} else if rowval >= 375 && colval >= 375 {
				grid_q4[i] = v
			}
		}
		wg.Add(1)
		go update(grid, &grid_q1, 0)
		wg.Add(1)
		go update(grid, &grid_q2, 1)
		wg.Add(1)
		go update(grid, &grid_q3, 2)
		wg.Add(1)
		go update(grid, &grid_q4, 3)
		wg.Wait()
		*grid = make(map[string]int)
		for i, v := range grid_q1 {
			(*grid)[i] = v
		}
		for i, v := range grid_q2 {
			(*grid)[i] = v
		}
		for i, v := range grid_q3 {
			(*grid)[i] = v
		}
		for i, v := range grid_q4 {
			(*grid)[i] = v
		}
	}

	f, err := os.OpenFile("seeds_of_chaos.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}

// chreate chaos pattern
func chaos(grid map[string]int, x int, y int) map[string]int {
	grid[strconv.Itoa(y)+","+strconv.Itoa(x+1)] = 1
	grid[strconv.Itoa(y+1)+","+strconv.Itoa(x+3)] = 1
	grid[strconv.Itoa(y+2)+","+strconv.Itoa(x)] = 1
	grid[strconv.Itoa(y+2)+","+strconv.Itoa(x+2)] = 1
	grid[strconv.Itoa(y+2)+","+strconv.Itoa(x+4)] = 1
	grid[strconv.Itoa(y+3)+","+strconv.Itoa(x)] = 1
	grid[strconv.Itoa(y+3)+","+strconv.Itoa(x+1)] = 1
	grid[strconv.Itoa(y+3)+","+strconv.Itoa(x+4)] = 1
	grid[strconv.Itoa(y+4)+","+strconv.Itoa(x+4)] = 1
	grid[strconv.Itoa(y+4)+","+strconv.Itoa(x+5)] = 1
	grid[strconv.Itoa(y+5)+","+strconv.Itoa(x+1)] = 1
	grid[strconv.Itoa(y+5)+","+strconv.Itoa(x+2)] = 1
	grid[strconv.Itoa(y+5)+","+strconv.Itoa(x+3)] = 1
	grid[strconv.Itoa(y+6)+","+strconv.Itoa(x+3)] = 1
	return grid
}

func main() {
	// create grid
	var grid = make(map[string]int)

	// add chaos pattern
	grid = chaos(grid, 375, 375)
	render(&grid, 300, 750)

}
