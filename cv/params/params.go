package params

import "image/color"

const (
	// Coloring determines how objects will be colored:
	/*
		0 - Random colors for each object
		1 - Recognized objects are colored based on their center color
	*/
	Coloring = 1
)

// Leniance determines how far off the size can be for each metric (height and width).
// For example an object 15x15 would be recognized for a RecognizedObject calling for 13x13 with a leniance setting of 2
var Leniance = 3

// FailedLimit is how many approximate frames must go by without GetWanted working before warning the user
// and using the unstuck algorithm
var FailedLimit = 100

// TileSize specifies how large each pathfinding tile should be in height and width.
// Larger tiles mean faster pathfinding calculation,
// however smaller tiles mean more precise and better pathfinding
var TileSize = 10

// TileColor is the color that the tiles should be rendered with
var TileColor = color.RGBA{0, 0, 0, 255}

// PathColor is the color of the tiles that Frisk will walk on
var PathColor = color.RGBA{0, 0, 255, 255}
