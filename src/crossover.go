package src

import "math/rand"

func Crossover(f1 *organism, f2 *organism) (*organism, *organism) {
	cut := rand.Intn(len(f1.DNA)-1) + 1
	X1 := f1.DNA[0:cut]
	Y1 := f1.DNA[cut:len(f1.DNA)]

	cut = rand.Intn(len(f2.DNA)-1) + 1
	X2 := f2.DNA[0:cut]
	Y2 := f2.DNA[cut:len(f2.DNA)]

	d1 := make([]int, len(X1)+len(Y2))
	i := 0
	i += copy(d1[i:], X1)
	copy(d1[i:], Y2)
	c1 := &organism{d1, 0}
	c1.removeDuplicateGenes()

	d2 := make([]int, len(X2)+len(Y1))
	i = 0
	i += copy(d2[i:], X2)
	copy(d2[i:], Y1)
	c2 := &organism{d2, 0}
	c2.removeDuplicateGenes()

	return c1, c2
}
