package mozaic

type TransformPipeline struct {
	Transforms []*Transform
}

func (p *TransformPipeline) Scale(base, values Point[float64]) {

}

func (p *TransformPipeline) Rotate(base Point[float64], degree float64) {

}

func (p *TransformPipeline) Translate(base, values Point[float64]) {

}

func (p *TransformPipeline) Reflect(reflector Line) {

}

func (p *TransformPipeline) ShearByLine(base Line) {

}

func (p *TransformPipeline) ShearByPoint(base Point[float64]) {

}

func (p *TransformPipeline) Clear() {
	p.Transforms = []*Transform{}
}
