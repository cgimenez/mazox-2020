package game

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"mazox/runtime"
	"reflect"
	"strings"
)

type ILevel interface {
	load()
	init()
	update()
	draw()

	addEntity(IEntity)
	entityCanMoveToPix(IEntity, int, int) bool
	entityCanMoveToCell(IEntity, int, int) bool
	getGridCellSize() int
	getGridSize() (int, int)
	translateGridLocToPix(int, int) (int, int)
	translatePixToGridLoc(int, int) (int, int)
}

type ogmoLevelLayer struct {
	Name           string     `json:"name"`
	Eid            string     `json:"_eid"`
	OffsetX        int        `json:"offsetX"`
	OffsetY        int        `json:"offsetY"`
	GridCellWidth  int        `json:"gridCellWidth"`
	GridCellHeight int        `json:"gridCellHeight"`
	GridCellsX     int        `json:"gridCellsX"`
	GridCellsY     int        `json:"gridCellsY"`
	Tileset        string     `json:"tileset"`
	Data2D         [][]int    `json:"data2D"`
	Grid2D         [][]string `json:"grid2D"`
	ExportMode     int        `json:"exportMode"`
	ArrayMode      int        `json:"arrayMode"`
}

type ogmoLevel struct {
	OgmoVersion string           `json:"ogmoVersion"`
	Width       int              `json:"width"`
	Height      int              `json:"height"`
	OffsetX     int              `json:"offsetX"`
	OffsetY     int              `json:"offsetY"`
	Layers      []ogmoLevelLayer `json:"layers"`
}

type Level struct {
	num                int
	gridCols, gridRows int
	gridCellSize       int
	grid               [10][20]int
	backGround         runtime.Image
	foreGround         runtime.Image
	debugLayer         runtime.Image
	zone               *Zone
	entities           []IEntity
}

const (
	LEVEL_GRID_FREE    = 0
	LEVEL_GRID_BLOCKED = 1
)

func levelFactory(num int) ILevel {
	var level ILevel
	switch num {
	case 1:
		level = &Level1{Level{num: 1}}
	}
	return level
}

//
// For a given tile ref, compute the drawing rect and pos inside destination image
//
func tileRefRectPoint(tile_ref int, col int, row int, tile_width int, tile_height int, grid_cell_width int, grid_cell_height int) (image.Rectangle, image.Point) {
	p := image.Point{
		tile_ref % tile_width * grid_cell_width,
		tile_ref / tile_height * grid_cell_width,
	}
	r := image.Rect(
		col*grid_cell_width,
		row*grid_cell_height,
		col*grid_cell_width+grid_cell_width,
		row*grid_cell_height+grid_cell_height,
	)
	return r, p
}

//
// Generate bitmap from tile data
//
func generateTileLayer(layer *ogmoLevelLayer, tile_filename string) runtime.Image {
	log.Printf("Loading layer [%s] - tile [%s]\n", layer.Name, layer.Tileset)
	destImg := image.NewRGBA(image.Rect(0, 0, gm.width, gm.height))
	tileImg := gm.loader.loadImage(tile_filename)
	tile_width := tileImg.Bounds().Max.X / layer.GridCellWidth
	tile_height := tileImg.Bounds().Max.Y / layer.GridCellHeight

	for row := 0; row < layer.GridCellsY; row++ {
		for col := 0; col < layer.GridCellsX; col++ {
			tile_ref := layer.Data2D[row][col]
			if tile_ref >= 0 {
				r, p := tileRefRectPoint(tile_ref, col, row, tile_width, tile_height, layer.GridCellWidth, layer.GridCellHeight)
				draw.Draw(destImg, r, tileImg, p, draw.Over)
			}
		}
	}

	return destImg
}

//
// Load JSON data
// initialize grid, background and foreground
//
func (l *Level) load() {
	l.gridCols = 20
	l.gridRows = 10
	l.gridCellSize = 32

	if gm.debug {
		l.debugLayer = image.NewRGBA(image.Rect(0, 0, gm.width, gm.height))
	}

	log.Printf("Loading level %d\n", l.num)
	jsonFile := gm.rtm.ReadFile(fmt.Sprintf("/levels/L%d.json", l.num))
	defer jsonFile.Close()

	var ogmoLevel ogmoLevel
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err := json.Unmarshal([]byte(byteValue), &ogmoLevel)
	if err != nil {
		fmt.Println(err)
	}

	for _, layer := range ogmoLevel.Layers {
		switch strings.Split(layer.Name, "-")[1] {
		case "Grid":
			for row := 0; row < l.gridRows; row++ {
				for col := 0; col < l.gridCols; col++ {
					if layer.Grid2D[row][col] == "#" {
						l.grid[row][col] = LEVEL_GRID_BLOCKED
						if gm.debug {
							blue := color.RGBA{0, 0, 255, 255}
							draw.Draw(l.debugLayer, image.Rect(col*32, row*32, col*32+32, row*32+32), image.NewUniform(blue), image.Point{0, 0}, draw.Over)
						}
					}
				}
			}
		case "Back":
			l.backGround = generateTileLayer(&layer, fmt.Sprintf("tiles/L%d.png", l.num))
		case "Fore":
			l.foreGround = generateTileLayer(&layer, fmt.Sprintf("tiles/L%d.png", l.num))
		}
	}
}

//
// Check if entity can move to X,Y pix coordinates
//
func (l *Level) entityCanMoveToPix(entity IEntity, x int, y int) bool {
	ew, eh := entity.getSize()
	if entity.getGridRespect() {
		if x < 0 || y < 0 || x >= l.gridCols*l.gridCellSize-ew || y >= l.gridRows*l.gridCellSize-eh {
			return false
		}
		ec, er := entity.getGridPosition()

		cells_to_check := surroundingCellsCoords(ec, er, l.gridCols, l.gridRows)
		rect_check := image.Rect(x, y, x+ew, y+eh)
		for _, cell := range cells_to_check {
			if cell.x >= 0 && cell.y >= 0 && cell.x < l.gridCols && cell.y < l.gridRows {
				/*
					fmt.Println(x, y, ec, er, rect_check,
						l.getGridCellRect(cell.x, cell.y),
						l.grid[cell.y][cell.x],
						rect_check.Overlaps(l.getGridCellRect(cell.x, cell.y)),
					)*/
				if rect_check.Overlaps(l.getGridCellRect(cell.x, cell.y)) && l.grid[cell.y][cell.x] == LEVEL_GRID_BLOCKED {
					return false
				}
			}
		}

		for i := range l.entities {
			if reflect.ValueOf(l.entities[i]).Pointer() != reflect.ValueOf(entity).Pointer() { // l.entities[i] != entity does not works
				if rect_check.Overlaps(l.entities[i].getRect()) {
					return false
				}
			}
		}
	}
	return true
}

//
// Check if entity can move to cell grid (path finding)
//
func (l *Level) entityCanMoveToCell(entity IEntity, col int, row int) bool {
	if l.grid[row][col] == LEVEL_GRID_BLOCKED {
		return false
	}
	for i := range l.entities {
		if reflect.ValueOf(l.entities[i]).Pointer() != reflect.ValueOf(entity).Pointer() {
			c, r := l.entities[i].getGridPosition()
			if c == col && r == row {
				return false
			}
		}
	}
	return true
}

//
// Returns the size of a grid cell
// Cells are square only
//
func (l *Level) getGridCellSize() int {
	return l.gridCellSize
}

//
// Returns size of grid in cols rows
//
func (l *Level) getGridSize() (int, int) {
	return l.gridCols, l.gridRows
}

//
// Returns the bounding rect of a grid cell
//
func (l *Level) getGridCellRect(col int, row int) image.Rectangle {
	return image.Rect(col*l.gridCellSize, row*l.gridCellSize, col*l.gridCellSize+l.gridCellSize, row*l.gridCellSize+l.gridCellSize)
}

//
// Compute the grid col row from X,Y pix
//
func (l *Level) translatePixToGridLoc(x int, y int) (int, int) {
	return x / l.gridCellSize, y / l.gridCellSize
}

//
// Compute the X,Y pix pos from grid col row
//
func (l *Level) translateGridLocToPix(col int, row int) (int, int) {
	return col * l.gridCellSize, row * l.gridCellSize
}

//
//
//
func (l *Level) addEntity(entity IEntity) {
	entity.setLevel(l)
	l.entities = append(l.entities, entity)
}

//
// Callbacks
//
func (l *Level) init() {
	SNBH()
}

func (l *Level) update() {
	SNBH()
}

func (l *Level) _update() {
	for entity := range l.entities {
		l.entities[entity].update()
	}
}

//
// Layers are
// 1) Background (draw src)
// 2) Entities (draw over)
// 3) Foreground (draw over)
//
func (l *Level) draw() {
	if l.backGround != nil {
		gm.rtm.DrawBitmap(l.backGround, 0, 0, true)
	}
	if gm.debug {
		gm.rtm.DrawBitmap(l.debugLayer, 0, 0, false)
	}
	for entity := range l.entities {
		l.entities[entity].draw()
	}
	if l.foreGround != nil {
		gm.rtm.DrawBitmap(l.foreGround, 0, 0, false)
	}
}
