package mozaic

type Loader interface {
	Load() error
}

type Saver interface {
	Save() error
}
