package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	q string
	a string
}

func parseLine(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return ret
}
func main() {
	// file := "/Users/likeminds/Desktop/code/go/quiz-cli/problems.csv"

	csvFilename := flag.String("csv", "/Users/likeminds/Desktop/code/go/quiz-cli/problems.csv", "A CSV file in format problem,answer")
	timeLimit := flag.Int("limit", 30, "Time Limit for the quiz")
	flag.Parse()

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	file, _ := os.Open(*csvFilename)
	defer file.Close()

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		fmt.Println("Error Parsing CSV")
	}

	problems := parseLine(lines)
	fmt.Println("Welcome to the Go quiz")
	var in string
	var score int = 0
	<-timer.C
	for i := 0; i < 10; i++ {

		fmt.Println(problems[i])
		fmt.Scanln(&in)

	}

	fmt.Println("Your score is:", score)

}
