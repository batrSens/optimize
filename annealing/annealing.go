package annealing

import (
	"optimize/function"
)

type Arsonist interface {
	Start(f *function.Func)
	NextTemp(f *function.Func, i int) (float64, bool)
	NewPoint(f *function.Func, xCur function.Point, temp float64) function.Point
	IsNext(f *function.Func, xCur, xNew function.Point, temp float64) bool
}

func MinFunc(a Arsonist) func(f *function.Func) function.Point {
	return func(f *function.Func) function.Point {
		return Minimum(a, f)
	}
}

func Minimum(a Arsonist, f *function.Func) function.Point {
	a.Start(f)

	xCur := f.Center()
	xMin := f.Center()

	t, isFinish := a.NextTemp(f, 0)
	for i := 1; !isFinish; i++ {
		if f.F(xCur) < f.F(xMin) {
			xMin = xCur
		}

		xNew := a.NewPoint(f, xCur, t)

		isNext := a.IsNext(f, xCur, xNew, t)
		t, isFinish = a.NextTemp(f, i)
		//if isFinish {
		//	fmt.Println(i)
		//}
		if isNext {
			xCur = xNew
		}
	}

	return xMin
}
