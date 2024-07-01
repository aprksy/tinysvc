package mozaic

type Base[T Number] interface {
	Attachables() []Attachable[T]
	Attachable() (Attachable[T], error)
	AddAttachable(att Attachable[T]) error
	RemoveAttachable(id string) error
	ClearAttachables()

	Attach(id string)
	Detach(id string)
	DetachAll()

	Transform(id string)
	TransformUsingCustomMatrix(id string, matrix [3][3]float64)
}

type Attachable[T Number] interface {
	Attach(basePoint, tessPoint Point[T])
	Detach()
}
