package function

type Point []float64

type Func struct {
	F func(x Point) float64

	From Point
	To   Point

	Size int
}

func NewFunc(f func(x Point) float64, size int, from, to float64) *Func {
	fromX := make(Point, size)
	toX := make(Point, size)

	for i := 0; i < size; i++ {
		fromX[i] = from
		toX[i] = to
	}

	return &Func{
		F:    f,
		From: fromX,
		To:   toX,
		Size: size,
	}
}

func NewFuncReverse(f *Func) *Func {
	return &Func{
		F: func(x Point) float64 {
			return -f.F(x)
		},
		From: f.From,
		To:   f.To,
		Size: f.Size,
	}
}

func NewFuncArea(f func(x Point) float64, from, to Point) *Func {
	return &Func{
		F:    f,
		From: from,
		To:   to,
		Size: len(from),
	}
}

func (f *Func) Center() Point {
	xs := make(Point, len(f.From))

	for i := 0; i < len(xs); i++ {
		xs[i] = (f.To[i]-f.From[i])/2 + f.From[i]
	}

	return xs
}

func (f *Func) IsBelongsCoord(val float64, i int) bool {
	return val >= f.From[i] && val <= f.To[i]
}

func (f *Func) IsBelongs(p Point) bool {
	for i := 0; i < len(p); i++ {
		if !f.IsBelongsCoord(p[i], i) {
			return false
		}
	}
	return true
}
