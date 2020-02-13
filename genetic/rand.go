package genetic

import (
	"math/rand"
)

type roulette struct {
	items   []int
	chances []float64
	r       *rand.Rand
}

func newRoulette(chances []float64, r *rand.Rand) *roulette {
	items := make([]int, len(chances))
	chancesCopy := make([]float64, len(chances))

	for i := range items {
		items[i] = i
		chancesCopy[i] = chances[i]
	}

	return &roulette{
		items:   items,
		chances: chancesCopy,
		r:       r,
	}
}

func (rl *roulette) get() int {
	res := randomPoint(rl.chances, rl.r)
	rl.chances = append(rl.chances[:res], rl.chances[res+1:]...)
	rl.items = append(rl.items[:res], rl.items[res+1:]...)
	return res
}

func randomPoint(fitness []float64, r *rand.Rand) int {
	sum := 0.0
	for _, v := range fitness {
		sum += v
	}

	num := r.Float64() * sum
	cur := 0.0

	for i := range fitness {
		cur += fitness[i]

		if num < cur {
			return i
		}
	}

	if len(fitness) == 0 {
		panic("null fitness")
	}
	panic("fitness sum less than 1.0")
}
