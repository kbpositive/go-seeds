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
func update(grid *map[string]int, grid_memo *map[string]int, bias *int, mx *int, slice int, quadrant int) {
	defer wg.Done()
	var memo = make(map[string]int)

	for key := range *grid_memo {
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
						stack = append(stack, child)
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
func render(grid *map[string]int, frames int, dim int, quadrants int) {
	var palette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0xff, 0xff},
		color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0xff, 0xff},
		color.RGBA{0xff, 0xff, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff},
	}
	var images []*image.Paletted
	var delays []int
	var mn = 375 - quadrants
	var mx = 375 + quadrants

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

		var grid_q []map[string]int
		for i := 0; i < quadrants; i++ {
			grid_q = append(grid_q, make(map[string]int))
		}

		var slice = mx / quadrants
		for i, v := range *grid {
			row := strings.Split(i, ",")[0]
			rowval, _ := strconv.Atoi(row)
			if rowval > mx {
				mx = rowval
			} else if rowval < mn {
				mn = rowval
			}
			for j := 0; j < quadrants; j++ {
				if rowval >= (j*slice)+mn && rowval < ((j+1)*slice)+mn {
					grid_q[j][i] = v
					break
				}
			}
		}
		wg.Add(quadrants)
		for i := 0; i < quadrants; i++ {
			go update(grid, &grid_q[i], &mn, &mx, slice, i)
		}
		wg.Wait()
		*grid = make(map[string]int)
		for j := 0; j < quadrants; j++ {
			for i, v := range grid_q[j] {
				(*grid)[i] = v
			}
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
	render(&grid, 300, 750, 16)

}
