package src

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type task struct {
	ID    int
	Value int
	Time  int
}

type Tasks []*task

func (tt *Tasks) GetByID(ID int) *task {
	return (*tt)[ID-1]
}

func (tt *Tasks) ReadCSV(f *os.File) {
	r := csv.NewReader(f)

	_, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal(err)
		}

		value, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal(err)
		}

		time, err := strconv.Atoi(line[2])
		if err != nil {
			log.Fatal(err)
		}

		*tt = append(*tt, &task{
			ID:    id,
			Value: value,
			Time:  time,
		})
	}
}
