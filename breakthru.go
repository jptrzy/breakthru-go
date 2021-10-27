package main

import (
	"fmt"
	"image/color"
	"jptrzy/breakthru/utils"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  int = 640
	screenHeight int = 640
	border_size  int = 6
	tile_size    int = 50
)

var board [11][11]int //[y][x]

var turn int = 0
var moves int = 1

var selected_x int = 0
var selected_y int = 0

var board_cell_colors = []color.RGBA{
	utils.ParseHexColorSimple("#1C3144"),
	utils.ParseHexColorSimple("#596F62"),
	utils.ParseHexColorSimple("#7EA16B"),
	utils.ParseHexColorSimple("#C3D898"),
	utils.ParseHexColorSimple("#16161E"),
	utils.ParseHexColorSimple("#e63946"),
	//{255, 255, 255, 255},
}

func setup_board() {
	for y := 0; y < 5; y++ {
		board[1][y+3] = 1
		board[9][y+3] = 1
		board[y+3][1] = 1
		board[y+3][9] = 1
	}

	for y := 0; y < 3; y++ {
		board[3][y+4] = 2
		board[7][y+4] = 2
		board[y+4][3] = 2
		board[y+4][7] = 2
	}

	board[5][5] = 3
}

func drawSelected(r *sdl.Renderer) {
	//if selected_x > -1 && selected_y > -1 && selected_x < 11 && selected_y < 11 {
	if moves%2 == 0 {
		return
	}

	col := board_cell_colors[5]
	r.SetDrawColor(col.R, col.G, col.B, col.A)
	r.FillRect(&sdl.Rect{int32(selected_x*(tile_size+border_size) + border_size/2), int32(selected_y*(tile_size+border_size) + border_size/2), int32(tile_size + border_size), int32(tile_size + border_size)})
}

func drawBackCell(r *sdl.Renderer, x int, y int) {
	col := board_cell_colors[4]
	r.SetDrawColor(col.R, col.G, col.B, col.A)

	r.FillRect(&sdl.Rect{int32(x*(tile_size+border_size) + border_size/2), int32(y*(tile_size+border_size) + border_size/2), int32(tile_size + border_size), int32(tile_size + border_size)})
}

func drawCell(r *sdl.Renderer, x int, y int) {
	col := board_cell_colors[board[y][x]]
	r.SetDrawColor(col.R, col.G, col.B, col.A)

	r.FillRect(&sdl.Rect{int32(x*(tile_size+border_size) + border_size), int32(y*(tile_size+border_size) + border_size), int32(tile_size), int32(tile_size)})
}

func drawAll(r *sdl.Renderer) {
	col := board_cell_colors[4]
	r.SetDrawColor(col.R, col.G, col.B, col.A)
	r.Clear()

	drawSelected(r)

	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board); x++ {
			drawCell(r, x, y)
		}
	}

	r.Present()
}

func main() {
	var print_neeaded bool = false

	setup_board()

	//SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL:", err)
		return
	}

	window, err := sdl.CreateWindow(
		"BreakThru",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(screenWidth), int32(screenHeight),
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("initializing window:", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("initializing renderer:", err)
		return
	}
	defer renderer.Destroy()

	drawAll(renderer)

	for {
		if print_neeaded {
			renderer.Present()
			print_neeaded = false
		}

		event := sdl.WaitEvent()
		switch ev := event.(type) {
		case *sdl.QuitEvent:
			return
		case *sdl.MouseButtonEvent:
			if ev.State == 1 && ev.Button == 1 {
				//TODO not enough accuracy
				drawBackCell(renderer, selected_x, selected_y)
				drawCell(renderer, selected_x, selected_y)
				selected_x = int(math.Floor(float64((int(ev.X) - border_size/2) / (tile_size + border_size))))
				selected_y = int(math.Floor(float64((int(ev.Y) - border_size/2) / (tile_size + border_size))))
				drawSelected(renderer)
				drawCell(renderer, selected_x, selected_y)
				print_neeaded = true
			}
		}
	}
}
