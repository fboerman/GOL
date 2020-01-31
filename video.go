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

// render given map to the screen
func render_map(Map *GOLMap, video *Video) {
	clear(video)

	var rects []sdl.Rect
	for y := 0; y < Map.heigth; y++ {
		for x := 0; x < Map.width; x++ {
			if *get_cell_read(x, y, Map) != 0 {
				rects = append(rects, sdl.Rect{
					X: int32(x * SCALE),
					Y: int32(y * SCALE),
					W: SCALE,
					H: SCALE,
				})
			}
		}
	}
	video.renderer.SetDrawColor(255, 255, 255, 255)
	if len(rects) != 0 {
		video.renderer.FillRects(rects)
	}
	video.renderer.Present()

}