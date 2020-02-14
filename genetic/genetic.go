package genetic

import (
	"optimize/function"
)

const (
	populationSize = 30
)

type Population interface {
	Start(f *function.Func)

	Parents(f *function.Func) Population
	Breeding(f *function.Func) Population
	Mutation(f *function.Func)
	Selection(f *function.Func, child Population)

	Population() []function.Point
	Fitness(f *function.Func) []float64

	IsFinish(f *function.Func) bool
	Result(f *function.Func) function.Point
}

func MinFunc(p Population) func(f *function.Func) function.Point {
	return func(f *function.Func) function.Point {
		return Minimum(p, f)
	}
}

func Minimum(p Population, f *function.Func) function.Point {
	p.Start(f)

	for !p.IsFinish(f) {
		parents := p.Parents(f)
		child := parents.Breeding(f)
		child.Mutation(f)
		p.Selection(f, child)
	}

	return p.Result(f)
}
