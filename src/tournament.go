package src

import "math/rand"

func Tournament(p *Population) *organism {
	i1 := rand.Intn(len(*p))
	i2 := i1
	for i2 == i1 {
		i2 = rand.Intn(len(*p))
	}
	o1 := (*p)[i1]
	o2 := (*p)[i2]
	if o1.Fitness > o2.Fitness {
		return o1
	} else {
		return o2
	}
}
