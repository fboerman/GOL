package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// parse a given file with format:
// first line: width height
// second line: x y start orientation
// lines of multiple characters with:
// X for living cell, any other single character for dead cell
func parse_map(fname string) *GOLMap {
	file, err := os.Open(fname)

	if err != nil {
		fmt.Println("[!] Could not open file!")
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	//allocate the board
	Map := new(GOLMap)

	//read the first two lines which denote width and height of the board
	if !scanner.Scan() {
		fmt.Println("[!] Invalid file format!")
		return nil
	}
	map_size := strings.Split(scanner.Text(), " ")
	var err1, err2 error
	Map.width, err1 = strconv.Atoi(map_size[0])
	Map.heigth, err2 = strconv.Atoi(map_size[1])

	if err1 != nil || err2 != nil {
		fmt.Println("[!] Invalid file format!")
		return nil
	}

	if !scanner.Scan() {
		fmt.Println("[!] Invalid file format!")
		return nil
	}
	orientation := strings.Split(scanner.Text(), " ")
	var x_start, y_start int
	x_start, err1 = strconv.Atoi(orientation[0])
	y_start, err2 = strconv.Atoi(orientation[1])

	if err1 != nil || err2 != nil {
		fmt.Println("[!] Invalid file format!")
		return nil
	}

	//allocate the memory for the board itself
	Map.data1 = make([]uint8, Map.width*Map.heigth)
	Map.data2 = make([]uint8, Map.width*Map.heigth)
	Map.current_buffer = 1
	x := 0
	y := 0
	//read the file and parse the board
	for scanner.Scan() {
		line := scanner.Text()
		for _, char := range line {
			if char == 'X' {
				*get_cell_read(x_start + x, y_start + y, Map) = 1
			}
			x++
		}
		x = 0
		y++
	}

	return Map
}
