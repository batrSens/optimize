package genetic

import (
	"math"
	"math/rand"
	"optimize/function"
	"time"
)

type Population1 struct {
	population []function.Point
	generation int
	maxGen     int

	min function.Point
	r   *rand.Rand
}

func NewPopulation1(size, generations int) *Population1 {
	return &Population1{
		population: make([]function.Point, size),
		maxGen:     generations,
		r:          rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func newPopulation1Pop(population []function.Point, r *rand.Rand) *Population1 {
	return &Population1{
		population: population,
		r:          r,
	}
}

func (p1 *Population1) Start(f *function.Func) {
	for i := range p1.population {
		p1.population[i] = f.RandomPoint(p1.r)
	}

	p1.min = nil
	p1.generation = 0
}

func (p1 *Population1) Parents(f *function.Func) Population {
	fitness := p1.Fitness(f)
	parents := make([]function.Point, len(p1.population)*2)

	for i := 0; i < len(p1.population); i++ {
		par1 := randomPoint(fitness, p1.r)
		par2 := randomPoint(fitness, p1.r)

		for par1 == par2 {
			//fmt.Println(p1.generation)
			par2 = randomPoint(fitness, p1.r)
		}

		parents[i*2] = p1.population[par1]
		parents[i*2+1] = p1.population[par2]
	}

	return newPopulation1Pop(parents, p1.r)
}

func (p1 *Population1) Breeding(f *function.Func) Population {
	const d = 0.25
	breed := make([]function.Point, len(p1.population)/2)

	for i := 0; i < len(p1.population)/2; i++ {
		par1 := p1.population[i*2]
		par2 := p1.population[i*2+1]
		child := make(function.Point, len(par1))

		for j := 0; j < len(par1); j++ {
			alpha := p1.r.Float64()*(1.0+d*2) - d
			child[j] = par1[j] + alpha*(par2[j]-par1[j])

			if child[j] < f.From[j] {
				child[j] = f.From[j]
			} else if child[j] > f.To[j] {
				child[j] = f.To[j]
			}
		}

		breed[i] = child
	}

	return newPopulation1Pop(breed, p1.r)
}

func (p1 *Population1) Mutation(f *function.Func) {
	const mutP = 0.01

	for i := 0; i < len(p1.population); i++ {
		if p1.r.Float64() < mutP {
			for j := 0; j < len(p1.population[i]); j++ {
				newGene := p1.population[i][j] + alphaMut(f, j)*deltaMut(f, p1.r)

				if newGene < f.From[j] {
					newGene = f.From[j]
				} else if newGene > f.To[j] {
					newGene = f.To[j]
				}

				p1.population[i][j] = newGene
			}
		}
	}
}

func (p1 *Population1) Selection(f *function.Func, child Population) {
	fullPopul := newPopulation1Pop(append(p1.population, child.Population()...), p1.r)
	fitness := fullPopul.Fitness(f)
	roul := newRoulette(fitness, p1.r)
	newPopul := make([]function.Point, len(p1.population))

	if f.F(fullPopul.min) < f.F(p1.min) {
		p1.min = fullPopul.min
	}

	for i := range p1.population {
		newPopul[i] = fullPopul.population[roul.get()]
	}

	p1.population = newPopul
	p1.generation++
}

func (p1 *Population1) IsFinish(f *function.Func) bool {
	return p1.generation > p1.maxGen
}

func (p1 *Population1) Result(f *function.Func) function.Point {
	return p1.min
}

func (p1 *Population1) Population() []function.Point {
	return p1.population
}

func (p1 *Population1) Fitness(f *function.Func) []float64 {
	res := make([]float64, len(p1.population))

	min := f.F(p1.population[0])
	minNum := 0
	for i := range p1.population {
		res[i] = f.F(p1.population[i])

		if res[i] < min {
			min = res[i]
			minNum = i
		}
	}

	if p1.min == nil || min < f.F(p1.min) {
		p1.min = p1.population[minNum]
	}

	sum := 0.0
	for i := range res {
		res[i] -= min
		sum += res[i]
	}

	if sum == 0.0 {
		p1.Cataclysm(f)
		return p1.Fitness(f)
	}

	for i := range res {
		res[i] = sum - res[i]
	}

	return res
}

func (p1 *Population1) Cataclysm(f *function.Func) {
	const share = 10
	size := len(p1.population) / share
	if size == 0 {
		size++
	}

	newPopul := append(p1.population[:size], make([]function.Point, len(p1.population)-size)...)
	for i := size; i < len(p1.population); i++ {
		newPopul[i] = f.RandomPoint(p1.r)
	}

	p1.population = newPopul
}

func alphaMut(f *function.Func, i int) float64 {
	return 0.5 * (f.To[i] - f.From[i])
}

func deltaMut(f *function.Func, r *rand.Rand) float64 {
	const m = 20.0

	res := 0.0
	for i := 0.0; i < m; i++ {
		if r.Float64() < (1.0 / m) {
			res += math.Pow(2, -i-1)
		}
	}

	if r.Float64() < 0.5 {
		return -res
	}

	return res
}
