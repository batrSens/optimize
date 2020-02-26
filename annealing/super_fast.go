package annealing

import (
	"math"
	"math/rand"
	"optimize/function"
	"time"
)

type SuperFast struct {
	tempStart  float64
	tempFinish float64
	decrement  float64
	r          *rand.Rand

	m, n float64
}

func NewSuperFast(tempStart, tempFinish, m, n float64) *SuperFast {
	return &SuperFast{
		tempStart:  tempStart,
		tempFinish: tempFinish,
		m:          m,
		n:          n,
		r:          rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (sf *SuperFast) Start(f *function.Func) {
	sf.decrement = sf.m * math.Exp(-sf.n/float64(f.Size))
}

func (sf *SuperFast) NextTemp(f *function.Func, i int) (float64, bool) {
	if i == 0 {
		return sf.tempStart, false
	}

	t := sf.tempStart * math.Exp(-sf.decrement*math.Pow(float64(i), 1.0/float64(f.Size)))
	return t, t < sf.tempFinish
}

func (sf *SuperFast) NewPoint(f *function.Func, xMin, xCur function.Point, temp float64) function.Point {
	xNew := make(function.Point, f.Size)
	rand.NormFloat64()
	for i := range xNew {
		for {
			r := sf.r.Float64()

			sign := 1.0
			if math.Signbit(r - 1/2) {
				sign = -1.0
			}

			xNew[i] = xCur[i] + sign*temp*(math.Pow(1.0+1.0/temp, 2.0*r-1.0)-1.0)

			if f.IsBelongsCoord(xNew[i], i) {
				break
			}
		}
	}
	return xNew
}

func (sf *SuperFast) IsNext(f *function.Func, xCur, xNew function.Point, temp float64) bool {
	p := math.Exp((f.F(xCur) - f.F(xNew)) / temp)
	r := sf.r.Float64()
	return r < p
}
