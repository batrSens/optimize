package annealing

import (
	"math"
	"math/rand"
	"optimize/function"
	"time"
)

type Boltzmann struct {
	tempStart  float64
	tempFinish float64
	r          *rand.Rand
}

func NewBoltzmann(tempStart, tempFinish float64) *Boltzmann {
	return &Boltzmann{
		tempStart:  tempStart,
		tempFinish: tempFinish,
		r:          rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (b *Boltzmann) NextTemp(f *function.Func, i int) (float64, bool) {
	if i == 0 {
		return b.tempStart, false
	}

	t := b.tempStart / math.Log(float64(1+i))
	return t, t < b.tempFinish
}

func (b *Boltzmann) NewPoint(f *function.Func, xCur function.Point, temp float64) function.Point {
	xNew := make(function.Point, f.Size)
	for i := range xNew {
		for {
			xNew[i] = b.r.NormFloat64()*temp + xCur[i]

			if f.IsBelongsCoord(xNew[i], i) {
				break
			}
		}
	}
	return xNew
}

func (b *Boltzmann) IsNext(f *function.Func, xCur, xNew function.Point, temp float64) bool {
	p := math.Exp((f.F(xCur) - f.F(xNew)) / temp)
	r := rand.Float64()
	return r < p
}
