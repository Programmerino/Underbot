package params

const (
	// Coloring determines how objects will be colored:
	/*
		0 - Random colors for each object
		1 - Recognized objects are colored accordingly to what is specified in objectAttrs.go, and unrecognized objects are gray
	*/
	Coloring = 1
)

// Leniance determines how far off the size can be for each metric (height and width). For example an object 15x15 would be recognized for a RecognizedObject calling for 13x13 with a leniance setting of 2
var Leniance = 3
