package pathfinding

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/beefsack/go-astar"
	"github.com/pkg/errors"
	"gitlab.com/256/Underbot/cv/params"
	"gitlab.com/256/Underbot/cv/rect"
)

// Tile is a section of the screen for pathfinding
type Tile struct {
	Pos       string      // String representation of the position. Ex: (1, 3)
	Coords    image.Point // The number representation of the above. This is not in pixels
	Rectangle image.Rectangle
	Color     color.Color // What color the tile should be rendered with
	Cost      int         // How long it takes to get past the tile
}

// How many rows and columns there are of tiles
var rows int
var columns int

var tiles []*Tile
var tileMap map[image.Point]*Tile

// MakeTiles creates a slice of tiles based on image dimensions
func MakeTiles(img image.RGBA) ([]*Tile, error) {
	if len(tiles) == 0 {
		// Needed for the amount of rows and columns needed
		width := img.Bounds().Dx()
		height := img.Bounds().Dy()
		rows = int(math.Ceil(float64(height / params.TileSize)))
		columns = int(math.Ceil(float64(width / params.TileSize)))

		for r := 0; r < rows; r++ {
			for c := 0; c < columns; c++ {
				tileRect := image.Rect(
					(c*params.TileSize)-params.TileSize,
					(r*params.TileSize)-params.TileSize,
					(c * params.TileSize),
					(r * params.TileSize))
				newTile := Tile{Pos: fmt.Sprintf("(%v, %v)", r, c), Rectangle: tileRect, Coords: image.Point{X: c, Y: r}, Color: params.TileColor}
				cost, err := newTile.GetCost(img)
				if err != nil {
					return []*Tile{}, errors.Wrap(err, "failed to calculate the cost")
				}
				newTile.Cost = cost
				tiles = append(tiles, &newTile)
			}
		}
	} else {
		for _, tile := range tiles {
			cost, err := tile.GetCost(img)
			if err != nil {
				return []*Tile{}, errors.Wrap(err, "failed to calculate the cost")
			}
			tile.Cost = cost
		}
	}
	if len(tileMap) == 0 {
		tileMap = MapTiles(tiles)
	}
	return tiles, nil
}

// GetCost calculates the movement cost for a tile
func (t *Tile) GetCost(screen image.RGBA) (int, error) {
	centerColor, err := rect.CenterColor(t.GetImage(screen))
	if err != nil {
		return 0, errors.Wrap(err, "could not get the center color")
	}

	// If the tile is black, then set cost to 255, because it's impossible to go there
	if (centerColor == color.RGBA{0, 0, 0, 255}) {
		t.Color = color.RGBA{255, 0, 0, 255}
		return 255, nil
	}
	t.Color = params.TileColor
	return 0, nil
}

// PathNeighbors returns a slice of the tiles neighboring a tile
func (t *Tile) PathNeighbors() []astar.Pather {
	tileX := t.Coords.X
	tileY := t.Coords.Y

	return []astar.Pather{
		tileMap[image.Point{tileX, tileY - 1}], // Up
		tileMap[image.Point{tileX + 1, tileY}], // Right
		tileMap[image.Point{tileX, tileY + 1}], // Down
		tileMap[image.Point{tileX - 1, tileY}], // Left
	}
}

// GetImage shows the image inside of a tile. The function should be passed the image of the whole screen
func (t *Tile) GetImage(screen image.RGBA) image.Image {
	return screen.SubImage(t.Rectangle)
}

func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	return 1
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	if t == nil {
		return 255
	}
	toT := to.(*Tile)
	absX := t.Coords.X - toT.Coords.X
	if absX < 0 {
		absX = -absX
	}
	absY := t.Coords.Y - toT.Coords.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}

// MapTiles allows the access to a tile based on the point used
func MapTiles(tiles []*Tile) map[image.Point]*Tile {
	tileMap := make(map[image.Point]*Tile)
	for _, tile := range tiles {
		tileMap[tile.Coords] = tile
	}
	return tileMap
}

// GetPath uses A* pathfinding to get the best route from one tile to another
func GetPath(a *Tile, b *Tile) ([]*Tile, error) {
	fmt.Println(a, "2:", b)
	path, _, found := astar.Path(a, b)
	if !found {
		return nil, errors.New("failed to find route")
	}
	var tempTiles []*Tile
	for _, pather := range path {
		tile := pather.(*Tile)
		tile.Color = params.PathColor
		tempTiles = append(tempTiles, tile)
	}
	return tempTiles, nil
}

// GetGoal returns the tile that the pathfinding should aim for
func GetGoal() *Tile {
	middlePoint := image.Point{columns - 1, (rows - 1) / 2}
	return tileMap[middlePoint] // The tile to the far right middle
}

// GetCurrentTile finds the tile that a point is in the most
func GetCurrentTile(fPoint image.Point) *Tile {

	var smallest *Tile
	var smallSize int

	// Calculate the smallest rectangle containing Frisk and the nearest tile
	for i, tile := range tiles {
		fRect := image.Rect(fPoint.X, fPoint.Y, fPoint.X+1, fPoint.Y+1)
		unionRect := tile.Rectangle.Union(fRect)
		averageSize := rect.AverageSize(unionRect)
		if i == 0 {
			smallest = tile
			smallSize = averageSize
		} else {
			if averageSize < smallSize {
				smallSize = averageSize
				smallest = tile
			}
		}
	}
	//fmt.Println(smallest)
	smallest.Color = color.RGBA{22, 193, 0, 255}
	return smallest
}
