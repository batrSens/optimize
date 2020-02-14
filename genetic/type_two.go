package genetic

import (
	"math/rand"
	"optimize/function"
	"time"
)

type Population2 struct {
	population []function.Point
	generation int
	maxGen     int

	min function.Point
	r   *rand.Rand
}

func NewPopulation2(size, generations int) *Population2 {
	return &Population2{
		population: make([]function.Point, size),
		maxGen:     generations,
		r:          rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func newPopulation2Pop(population []function.Point, r *rand.Rand) *Population2 {
	return &Population2{
		population: population,
		r:          r,
	}
}

func (p2 *Population2) Start(f *function.Func) {
	for i := range p2.population {
		p2.population[i] = f.RandomPoint(p2.r)
	}

	p2.min = nil
	p2.generation = 0
}

func (p2 *Population2) Parents(f *function.Func) Population {
	fitness := p2.Fitness(f)
	parents := make([]function.Point, len(p2.population)*2)

	for i := 0; i < len(p2.population); i++ {
		par1 := randomPoint(fitness, p2.r)
		par2 := randomPoint(fitness, p2.r)

		for par1 == par2 {
			//fmt.Println(p2.generation)
			par2 = randomPoint(fitness, p2.r)
		}

		parents[i*2] = p2.population[par1]
		parents[i*2+1] = p2.population[par2]
	}

	return newPopulation1Pop(parents, p2.r)
}

func (p2 *Population2) Breeding(f *function.Func) Population {
	const d = 0.25
	breed := make([]function.Point, len(p2.population)/2)

	for i := 0; i < len(p2.population)/2; i++ {
		par1 := p2.population[i*2]
		par2 := p2.population[i*2+1]
		child := make(function.Point, len(par1))

		alpha := p2.r.Float64()*(1.0+d*2) - d
		for j := 0; j < len(par1); j++ {
			child[j] = par1[j] + alpha*(par2[j]-par1[j])

			if child[j] < f.From[j] {
				child[j] = f.From[j]
			} else if child[j] > f.To[j] {
				child[j] = f.To[j]
			}
		}

		breed[i] = child
	}

	return newPopulation1Pop(breed, p2.r)
}

func (p2 *Population2) Mutation(f *function.Func) {
	const mutP = 0.01

	for i := 0; i < len(p2.population); i++ {
		if p2.r.Float64() < mutP {
			for j := 0; j < len(p2.population[i]); j++ {
				newGene := p2.population[i][j] + alphaMut(f, j)*deltaMut(f, p2.r)

				if newGene < f.From[j] {
					newGene = f.From[j]
				} else if newGene > f.To[j] {
					newGene = f.To[j]
				}

				p2.population[i][j] = newGene
			}
		}
	}
}

func (p2 *Population2) Selection(f *function.Func, child Population) {
	fullPopul := newPopulation1Pop(append(p2.population, child.Population()...), p2.r)
	fitness := fullPopul.Fitness(f)
	roul := newRoulette(fitness, p2.r)
	newPopul := make([]function.Point, len(p2.population))

	if f.F(fullPopul.min) < f.F(p2.min) {
		p2.min = fullPopul.min
	}

	for i := range p2.population {
		newPopul[i] = fullPopul.population[roul.get()]
	}

	p2.population = newPopul
	p2.generation++
}

func (p2 *Population2) IsFinish(f *function.Func) bool {
	return p2.generation > p2.maxGen
}

func (p2 *Population2) Result(f *function.Func) function.Point {
	return p2.min
}

func (p2 *Population2) Population() []function.Point {
	return p2.population
}

func (p2 *Population2) Fitness(f *function.Func) []float64 {
	res := make([]float64, len(p2.population))

	min := f.F(p2.population[0])
	minNum := 0
	for i := range p2.population {
		res[i] = f.F(p2.population[i])

		if res[i] < min {
			min = res[i]
			minNum = i
		}
	}

	if p2.min == nil || min < f.F(p2.min) {
		p2.min = p2.population[minNum]
	}

	sum := 0.0
	for i := range res {
		res[i] -= min
		sum += res[i]
	}

	if sum == 0.0 {
		p2.Cataclysm(f)
		return p2.Fitness(f)
	}

	for i := range res {
		res[i] = sum - res[i]
	}

	return res
}

func (p2 *Population2) Cataclysm(f *function.Func) {
	const share = 10
	size := len(p2.population) / share
	if size == 0 {
		size++
	}

	newPopul := append(p2.population[:size], make([]function.Point, len(p2.population)-size)...)
	for i := size; i < len(p2.population); i++ {
		newPopul[i] = f.RandomPoint(p2.r)
	}

	p2.population = newPopul
}
