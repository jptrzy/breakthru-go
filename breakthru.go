package main

import (
	"fmt"
	"image/color"
	"jptrzy/breakthru/utils"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	board_size = 11
)

var board [board_size * board_size]int

var player int = 0
var moves int = 0

var selected_x int = 0
var selected_y int = 0

//Graphics

const (
	FPS                  = 40
	block_proc_in_board  = 12
	border_proc_in_block = 13
)

var wait_time int32 = int32(float32(1000) / FPS)
var frame_start_time uint32 = 0
var delay_time int32

var screen_width int32 = 640
var screen_height int32 = 640
var screen_size int32
var block_size int32
var border_size int32
var board_g_size int32
var margin_x int32
var margin_y int32

var color_sheme = []color.RGBA{
	utils.ParseHexColorSimple("#1C3144"),
	utils.ParseHexColorSimple("#596F62"),
	utils.ParseHexColorSimple("#7EA16B"),
	utils.ParseHexColorSimple("#C3D898"),
	utils.ParseHexColorSimple("#16161E"),
	utils.ParseHexColorSimple("#e63946"),
}

func setup_board() {
	// for y := 0; y < 5; y++ {
	// 	board[1][y+3] = 1
	// 	board[9][y+3] = 1
	// 	board[y+3][1] = 1
	// 	board[y+3][9] = 1
	// }

	// for y := 0; y < 3; y++ {
	// 	board[3][y+4] = 2
	// 	board[7][y+4] = 2
	// 	board[y+4][3] = 2
	// 	board[y+4][7] = 2
	// }

	board[60] = 3
}

func SetDrawColorByRGBA(rend *sdl.Renderer, col color.RGBA) {
	rend.SetDrawColor(col.R, col.G, col.B, col.A)
}

func updateSize() {
	screen_size = screen_width
	if screen_width > screen_height {
		screen_size = screen_height
	}

	block_size = screen_size / block_proc_in_board
	border_size = block_size / border_proc_in_block
	board_g_size = (block_size*11 + border_size*12)
	margin_x = (screen_width - board_g_size) / 2
	margin_y = (screen_height - board_g_size) / 2
}

func drawBoard(rend *sdl.Renderer) {
	SetDrawColorByRGBA(rend, color_sheme[0])
	rend.Clear()

	//fmt.Print(margin_x, screen_size-board_g_size)
	SetDrawColorByRGBA(rend, color_sheme[4])
	rend.FillRect(&sdl.Rect{margin_x, margin_y, board_g_size, board_g_size})

	for y := 0; y < 11; y++ {
		for x := 0; x < 11; x++ {
			SetDrawColorByRGBA(rend, color_sheme[board[11*y+x]])
			rend.FillRect(&sdl.Rect{margin_x + border_size + (block_size+border_size)*int32(x), margin_y + border_size + (block_size+border_size)*int32(y), block_size, block_size})
		}
	}

	rend.Present()
}

func main() {

	//Varibles
	var running bool = true

	updateSize()
	setup_board()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Breakthru", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(screen_width), int32(screen_height), sdl.WINDOW_VULKAN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	drawBoard(renderer)

	//renderer.FillRect(nil)

	// surface.FillRect(nil, 0)

	// rect := sdl.Rect{0, 0, 200, 200}
	// surface.FillRect(&rect, 0xffff0000)
	// window.UpdateSurface()
	// rect = sdl.Rect{50, 100, 200, 200}
	// surface.FillRect(&rect, 0xffffff00)
	// window.UpdateSurface()

	//drawBoard(surface);
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.WindowEvent:
				if ev.Event == sdl.WINDOWEVENT_RESIZED {
					fmt.Println("RESIZE")
					screen_width, screen_height = window.GetSize()
					updateSize()
					drawBoard(renderer)

				}
			case *sdl.MouseButtonEvent:
				if ev.State == 1 && ev.Button == 1 {

				}
			}

		}

		//Sleep to don't go in loop for nothing to often. Just optimize loop/rendering.
		delay_time = wait_time - int32(sdl.GetTicks()-frame_start_time)
		if delay_time > 0 {
			sdl.Delay(uint32(delay_time))
		}
		frame_start_time = sdl.GetTicks()
	}
}
