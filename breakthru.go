package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"jptrzy/breakthru/utils"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	board_size = 11
)

var board [board_size * board_size]int

var player int = 1

var moved_pawn int = -1

var selected bool = false
var selected_tile int = 0

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

type Config struct {
	version     int
	color_sheme []string
}

func load_json() {
	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Now let's unmarshall the data into `payload`
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	e, c := payload["color_sheme"].([]interface{})
	// log.Printf("%i\n", payload.version)
	// log.Println(string(content))
	log.Println(payload["color_sheme"])
	log.Println(e, c)
	// log.Println(payload.version)
	for i := 0; i < len(e); i++ {
		color_sheme[i] = utils.ParseHexColorSimple(e[i].(string))
	}

	// Let's print the unmarshalled data!
	//log.Printf("color_sheme: %s\n", payload["color_sheme"])
}

func setup_board() {
	for i := 0; i < 5; i++ {
		board[14+i] = 1
		board[102+i] = 1
		board[11*(i+3)+1] = 1
		board[11*(i+3)+9] = 1
	}

	for i := 0; i < 3; i++ {
		board[37+i] = 2
		board[81+i] = 2
		board[11*(i+4)+3] = 2
		board[11*(i+4)+7] = 2
	}

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

func drawTile(rend *sdl.Renderer, tile int, back int) {
	x := tile % 11
	y := (tile - x) / 11
	fmt.Println(x, y)
	SetDrawColorByRGBA(rend, color_sheme[back])
	rend.FillRect(&sdl.Rect{margin_x + border_size/2 + (block_size+border_size)*int32(x), margin_y + border_size/2 + (block_size+border_size)*int32(y), block_size + border_size, block_size + border_size})
	SetDrawColorByRGBA(rend, color_sheme[board[tile]])
	rend.FillRect(&sdl.Rect{margin_x + border_size + (block_size+border_size)*int32(x), margin_y + border_size + (block_size+border_size)*int32(y), block_size, block_size})
}

//TODO find easier way
func drawBack(rend *sdl.Renderer) {
	SetDrawColorByRGBA(rend, color_sheme[player+1])
	rend.FillRect(&sdl.Rect{0, 0, screen_width, margin_y})
	rend.FillRect(&sdl.Rect{0, margin_y + board_g_size, screen_width, margin_y})

	rend.FillRect(&sdl.Rect{0, margin_y, margin_x, board_g_size})
	rend.FillRect(&sdl.Rect{margin_x + board_g_size, margin_y, margin_x, board_g_size})
}

func drawBoard(rend *sdl.Renderer) {
	SetDrawColorByRGBA(rend, color_sheme[player+1])
	rend.Clear()

	SetDrawColorByRGBA(rend, color_sheme[4])
	rend.FillRect(&sdl.Rect{margin_x, margin_y, board_g_size, board_g_size})

	for y := 0; y < 11; y++ {
		for x := 0; x < 11; x++ {
			SetDrawColorByRGBA(rend, color_sheme[board[11*y+x]])
			rend.FillRect(&sdl.Rect{margin_x + border_size + (block_size+border_size)*int32(x), margin_y + border_size + (block_size+border_size)*int32(y), block_size, block_size})
		}
	}

	if selected {
		drawTile(rend, selected_tile, 5)
	}

	rend.Present()
}

func on_mouse_click(rend *sdl.Renderer, x int32, y int32) {
	if x > margin_x && x < margin_x+board_g_size && y > margin_y && y < margin_y+board_g_size {
		x -= margin_x
		y -= margin_y
		tile_x := x / (border_size + block_size)
		tile_y := y / (border_size + block_size)
		on_tile_click(rend, int(tile_y*11+tile_x))
	}
}

//Check if tile is actual player pawn
func is_allay(tile int) bool {
	if player == 0 && board[tile] == 1 || player == 1 && board[tile] > 1 {
		return true
	}
	return false
}

func on_tile_click(rend *sdl.Renderer, tile int) {
	drawTile(rend, selected_tile, 4)
	if selected {
		if tile != selected_tile {
			//START Moving Logick

			if board[tile] == 0 && (selected_tile == tile-1 || selected_tile == tile+1 || selected_tile == tile-11 || selected_tile == tile+11) {
				board[tile] = board[selected_tile]
				board[selected_tile] = 0
				if moved_pawn != -1 {
					moved_pawn = -1
					player = (player + 1) % 2
					drawBack(rend)
				} else {
					moved_pawn = tile
				}
			} else if board[tile] != 0 && !is_allay(tile) && (selected_tile == tile-10 || selected_tile == tile-12 || selected_tile == tile+10 || selected_tile == tile+12) {
				board[tile] = board[selected_tile]
				board[selected_tile] = 0
				moved_pawn = -1
				player = (player + 1) % 2
				drawBack(rend)
			}

			//END

			drawTile(rend, selected_tile, 4)
			drawTile(rend, tile, 4)
		}
		selected = false
	} else if board[tile] != 0 && is_allay(tile) && tile != moved_pawn {
		selected_tile = tile
		drawTile(rend, tile, 5)
		selected = true
	}
	rend.Present()
}

func main() {
	load_json()

	//Varibles
	var running bool = true

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

	updateSize()
	drawBoard(renderer)

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.WindowEvent:
				if ev.Event == sdl.WINDOWEVENT_RESIZED || ev.Event == sdl.WINDOWEVENT_EXPOSED {
					screen_width, screen_height = window.GetSize()
					updateSize()
					drawBoard(renderer)
				}
			case *sdl.MouseButtonEvent:
				if ev.State == 1 && ev.Button == 1 {
					on_mouse_click(renderer, ev.X, ev.Y)
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
