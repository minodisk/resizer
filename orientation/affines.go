package orientation

import (
	"math"

	"github.com/BurntSushi/graphics-go/graphics"
)

func toRadian(d float64) float64 {
	return math.Pi * d / 180
}

var affines map[int]graphics.Affine = map[int]graphics.Affine{
	1: graphics.I,
	2: graphics.I.Scale(-1, 1),
	3: graphics.I.Scale(-1, -1),
	4: graphics.I.Scale(1, -1),
	5: graphics.I.Rotate(toRadian(90)).Scale(-1, 1),
	6: graphics.I.Rotate(toRadian(90)),
	7: graphics.I.Rotate(toRadian(-90)).Scale(-1, 1),
	8: graphics.I.Rotate(toRadian(-90)),
}

// var affines map[int]graphics.Affine = map[int]graphics.Affine{
// 	1: graphics.Affine{1, 0, 0, 0, 1, 0, 0, 0, 1},
// 	2: graphics.Affine{-1, 0, 0, 0, 1, 0, 0, 0, 1},
// 	3: graphics.Affine{-1, 0, 0, 0, -1, 0, 0, 0, 1},
// 	4: graphics.Affine{1, 0, 0, 0, -1, 0, 0, 0, 1},
// 	5: graphics.Affine{0, 1, 0, 1, 0, 0, 0, 0, 1},
// 	6: graphics.Affine{0, 1, 0, -1, 0, 0, 0, 0, 1},
// 	7: graphics.Affine{0, -1, 0, -1, 0, 0, 0, 0, 1},
// 	8: graphics.Affine{0, -1, 0, 1, 0, 0, 0, 0, 1},
// }

// func init() {
// 	fmt.Println("-------------------")
// 	fmt.Printf("%+v\n", affines)
// 	fmt.Println("-------------------")
// }
