package main

import "github.com/veandco/go-sdl2/sdl"

type Video struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	tex      *sdl.Texture
	Dirty    bool
}

const SCALE = 10
const WIDTH = 100
const HEIGTH = 100

// initializes the video driver
func init_video() *Video {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Conway Game of Life", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		SCALE*WIDTH, SCALE*HEIGTH, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_STREAMING, int32(SCALE*WIDTH), int32(SCALE*HEIGTH))
	if err != nil {
		panic(err)
	}

	video := new(Video)
	video.renderer.SetDrawColor(0, 0, 0, 0)
	video.window = window
	video.renderer = renderer
	video.tex = tex

	return video
}

// neatly close off SDL
func close_video(video *Video) {
	video.tex.Destroy()
	video.renderer.Destroy()
	video.window.Destroy()
	sdl.Quit()
}

// clear the screen
func clear(video *Video) {
	video.renderer.SetDrawColor(0, 0, 0, 0)
	video.renderer.Clear()
	video.renderer.Present()
}

// color palete
var colors = [...][3]uint8{
	{255, 77, 77},
	{255, 26, 26},
	{230, 0, 0},
	{179, 0, 0},
	{128, 0, 0},
	{77, 0, 0},
}

// render given map to the screen
func render_map(Map *GOLMap, video *Video) {
	clear(video)

	for y := 0; y < Map.heigth; y++ {
		for x := 0; x < Map.width; x++ {
			current_cell := *get_cell_read(x, y, Map)
			if current_cell != 0 {
				rect := sdl.Rect{
					X: int32(x*SCALE) + 1,
					Y: int32(y*SCALE) + 1,
					W: SCALE - 1,
					H: SCALE - 1,
				}
				var color [3]uint8
				if int(current_cell) >= len(colors) {
					color = colors[len(colors)-1]
				} else {
					color = colors[current_cell-1]
				}
				video.renderer.SetDrawColor(color[0], color[1], color[2], 255)
				video.renderer.FillRect(&rect)
			}
		}
	}
	video.renderer.Present()

}
