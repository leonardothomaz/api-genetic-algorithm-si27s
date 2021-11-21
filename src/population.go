package src

import (
	"math/rand"
)

type Population []*organism

func (p *Population) GeneratePopulation(size int, min int, max int) {
	for i := 0; i < size; i++ {
		var o organism
		o.generateDNA(min, max)
		*p = append(*p, &o)
	}
}

func (p *Population) CalculateFitness(tt *Tasks) {
	for _, o := range *p {
		o.calculateFitness(tt)
	}
}

func (p *Population) GetBestAndWorst() (*organism, *organism) {
	best := (*p)[0]
	worst := (*p)[0]
	for i := 1; i < len(*p); i++ {
		if (*p)[i].Fitness > best.Fitness {
			best = (*p)[i]
		}
		if (*p)[i].Fitness < worst.Fitness {
			worst = (*p)[i]
		}
	}
	return best, worst
}

func (p *Population) Mutate(min int, max int, rate float64) {
	best, worst := p.GetBestAndWorst()
	for i := 0; i < len(*p); i++ {
		m := rand.Intn(100)
		if (m < int(rate*100) && (*p)[i] != best) || (*p)[i] == worst {
			(*p)[i].mutate(min, max)
		}
	}
}
