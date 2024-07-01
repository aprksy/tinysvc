package mozaic

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type Point[T Number] struct {
	X, Y T
}

type Line [3]float64

type Transform [3][3]float64
