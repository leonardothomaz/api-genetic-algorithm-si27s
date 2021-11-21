package src

import "math/rand"

type organism struct {
	DNA     []int
	Fitness float64
}

func (o *organism) generateDNA(min int, max int) {
	size := rand.Intn(max) + min
	for len(o.DNA) < size {
		gene := rand.Intn(max) + min
		duplicate := false
		for _, g := range o.DNA {
			if g == gene {
				duplicate = true
				break
			}
		}
		if !duplicate {
			o.DNA = append(o.DNA, gene)
		}
	}
}

func (o *organism) calculateFitness(tt *Tasks) {
	o.Fitness = 0
	time := 0
	for _, g := range o.DNA {
		o.Fitness += float64(tt.GetByID(g).Value)
		time += tt.GetByID(g).Time
	}
	if time <= 85 || time >= 90 {
		if time > 75 && time < 100 {
			o.Fitness /= 2
		} else {
			o.Fitness /= 10
		}
	}
}

func (o *organism) mutate(min int, max int) {
	m := rand.Intn(max) + min
	pos := rand.Intn(len(o.DNA))
	o.DNA[pos] = m
	o.removeDuplicateGenes()
}

func (o *organism) removeDuplicateGenes() {
	keys := make(map[int]bool)
	var dna []int
	for _, entry := range o.DNA {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			dna = append(dna, entry)
		}
	}
	o.DNA = dna
}
