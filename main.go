package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"

	. "genetic-algorithm-si27s/src"
)

const (
	PopulationSize = 32
	Generations    = 100_000
	MutationRate   = 0.03
)

func main() {
	dataset, err := os.Open("data/dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dataset)

	resultFile, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(w *os.File) {
		err := w.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resultFile)

	err = resultFile.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, resultFile)

	write(mw,
		"> START\n"+
			"- - - - - - - - - - - - - - - - - - - -\n"+
			" - PARAMETERS:\n"+
			"   - population size: %v\n"+
			"   - generations: %d\n"+
			"   - mutation rate: %.2f\n"+
			"- - - - - - - - - - - - - - - - - - - -\n",
		PopulationSize, Generations, MutationRate)

	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	var elapsed int64

	var tt Tasks
	tt.ReadCSV(dataset)

	var p Population
	p.GeneratePopulation(PopulationSize, 1, 99)
	p.CalculateFitness(&tt)

	generation := 0
	solution, _ := p.GetBestAndWorst()
	foundAt := generation
	write(mw, "gen: %d | best: %.2f\n", generation, solution.Fitness)

	for generation < Generations {
		var pool Population

		for i := 0; i < len(p)/2; i++ {
			f1 := Tournament(&p)
			f2 := Tournament(&p)
			c1, c2 := Crossover(f1, f2)
			pool = append(pool, c1)
			pool = append(pool, c2)
		}

		p = pool
		p.CalculateFitness(&tt)
		p.Mutate(1, 99, MutationRate)

		generation++

		best, worst := p.GetBestAndWorst()
		if best.Fitness > solution.Fitness {
			solution = best
			foundAt = generation
			elapsed = time.Since(start).Milliseconds()
			write(mw, "gen: %d | best: %.2f\n", generation, solution.Fitness)
		} else {
			*worst = *solution
		}
	}

	sort.Ints(solution.DNA)

	write(mw,
		"- - - - - - - - - - - - - - - - - - - -\n"+
			"  - SOLUTION:\n"+
			"    - dna: %v\n"+
			"    - fitness: %.2f\n"+
			"    - generation: %d\n"+
			"    - found in: %dms\n"+
			"- - - - - - - - - - - - - - - - - - - -\n"+
			"> END\n",
		solution.DNA, solution.Fitness, foundAt, elapsed)
}

func write(w io.Writer, str string, args ...interface{}) {
	_, err := fmt.Fprintf(w, str, args...)
	if err != nil {
		log.Fatal(err)
	}
}
