package main

import (
	"flag"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"time"
)

type GOLMap struct {
	width          int
	heigth         int
	data1          []uint8
	data2          []uint8
	current_buffer int
}

// return pointer to cell in the map with given coordinates
func get_cell_read(x int, y int, Map *GOLMap) *uint8 {
	var data []uint8
	if Map.current_buffer == 1 {
		data = Map.data1
	} else if Map.current_buffer == 2 {
		data = Map.data2
	} else {
		return nil
	}

	//return &(data[(y%Map.heigth)*Map.heigth+(x%Map.width)])
	return &(data[y*Map.heigth+x])
}

// return pointer to the cell in the write buffer
func get_cell_write(x int, y int, Map *GOLMap) *uint8 {
	var data []uint8
	if Map.current_buffer == 1 {
		data = Map.data2
	} else if Map.current_buffer == 2 {
		data = Map.data1
	} else {
		return nil
	}

	return &(data[y*Map.heigth+x])
}

// return currently selected buffer
func get_current_buffer(Map *GOLMap) *[]uint8 {
	if Map.current_buffer == 1 {
		return &Map.data1
	} else if Map.current_buffer == 2 {
		return &Map.data2
	}
	return nil
}

// this uses bit trick to prevent if statement if number is larger then zero
// https://stackoverflow.com/questions/3912112/check-if-a-number-is-non-zero-using-bitwise-operators-in-c
func sum(V []uint8) (sum uint8) {
	for _, x := range V {
		sum += (x | (uint8(^x) + 1)) >> 7
	}
	return
}

// calculate new generation
// from wiki:
//   Any live cell with two or three neighbors survives.
//   Any dead cell with three live neighbors becomes a live cell.
//   All other live cells die in the next generation. Similarly, all other dead cells stay dead.
func next(Map *GOLMap) {
	// iterators
	// 3x3 block needed, a,b,c are the start indices per row in the block
	it_a := 0
	it_b := Map.width
	it_c := 2 * Map.width

	//current data buffer
	data := *get_current_buffer(Map)

	// y and x are the middle pixel of the block
	for y := 1; y < Map.heigth-2; y++ {
		for x := 1; x < Map.width-2; x++ {

			// sum all alive cells in the block of 3x3
			sum_block := sum(data[it_a:it_a+3]) + sum(data[it_b:it_b+3]) + sum(data[it_c:it_c+3])

			// from wiki
			// To avoid decisions and branches in the counting loop, the rules can be rearranged from an
			// egocentric approach of the inner field regarding its neighbours to a scientific observer's
			// viewpoint: if the sum of all nine fields in a given neighbourhood is three, the inner field state for
			// the next generation will be life; if the all-field sum is four, the inner field retains its current state;
			// and every other sum sets the inner field to death.
			current_cell := *get_cell_read(x, y, Map)
			switch sum_block {
			case 3:
				*get_cell_write(x, y, Map) += 1
			case 4:
				if current_cell > 0 {
					*get_cell_write(x, y, Map) = current_cell + 1
				} else {
					*get_cell_write(x, y, Map) = 0
				}
			default:
				*get_cell_write(x, y, Map) = 0
			}

			it_a++
			it_b++
			it_c++
		}

		it_a += 3
		it_b += 3
		it_c += 3

	}

	//swap the buffers
	switch Map.current_buffer {
	case 1:
		Map.current_buffer = 2
	case 2:
		Map.current_buffer = 1
	}
}

func main() {
	map_fname := flag.String("map", "", "name of file to load map from")

	flag.Parse()

	if _, err := os.Stat(*map_fname); os.IsNotExist(err) {
		fmt.Println("[!] Invalid map file!")
		return
	}

	fmt.Println("[>] Starting Conway Game of Life")

	fmt.Println("[>] Loading map file")
	Map := parse_map(*map_fname)
	fmt.Println("[>] map loaded of ", Map.width, " by ", Map.heigth)

	fmt.Println("[>] Starting SDL")
	video := init_video()
	defer close_video(video)
	fmt.Println("[>] SDL loaded")

	i := 0

	for true {
		fmt.Println("[>] generation: ", i)
		// render the map
		render_map(Map, video)
		time.Sleep(500 * time.Millisecond)

		// check if button has been preseed
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			if event.GetType() == sdl.KEYDOWN || event.GetType() == sdl.KEYUP {
				return
			}
		}

		// calculate the next generation
		next(Map)
		i++
	}

}
