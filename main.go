package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	. "genetic-algorithm-si27s/src"

	"github.com/gorilla/mux"
)

const (
	PopulationSize = 32
	Generations    = 100_000
	MutationRate   = 0.03
)

type Solution struct {
	Dna        string
	Fitness    string
	Generation string
	Found_in   string
}

func main() {
	rotes := mux.NewRouter().StrictSlash(true)

	rotes.HandleFunc("/", getSolution).Methods("GET")
	var port = ":8000"
	fmt.Println("Server running in port:", port)
	log.Fatal(http.ListenAndServe(port, rotes))

}

func handleSolution() Solution {
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

	dna := arrayToString(solution.DNA, " ")
	fitness := fmt.Sprintf("%f", solution.Fitness)
	generationSol := fmt.Sprintf("%d", foundAt)
	found_in := fmt.Sprintf("%dms", elapsed)

	var solucao = Solution{
		Dna:        dna,
		Fitness:    fitness,
		Generation: generationSol,
		Found_in:   found_in,
	}

	return solucao
}

func getSolution(w http.ResponseWriter, r *http.Request) {
	solution := handleSolution()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(solution)
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

func write(w io.Writer, str string, args ...interface{}) {
	_, err := fmt.Fprintf(w, str, args...)
	if err != nil {
		log.Fatal(err)
	}
}
