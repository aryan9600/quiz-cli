package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	csvFile := flag.String("csv", "problems.csv", "csv file for the problems")
	timeLimit := flag.Int("limit", 60, "time limit for the quiz")
	flag.Parse()

	file, err := os.Open(*csvFile)
	if err != nil {
		exit(fmt.Sprintf("Unable to open the csv file %s\n", *csvFile))
		os.Exit(1)
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the given csv file.")
	}
	problems := parseLines(lines)
	playQuiz(problems, 0, timeLimit)

}

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return ret
}

func playQuiz(questions []problem, count int, limit *int) {

	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	for i, p := range questions {
		fmt.Printf("Problem %d: %s = ", i+1, p.q)

		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTime's up. You scored %d out of %d.", count, len(questions))
			return
		case answer := <-answerCh:
			if answer == p.a {
				count++
				fmt.Println("Correct!")
			} else {
				fmt.Println("Incorrect!")
			}
		}
	}

	fmt.Printf("\nYou scored %d out of %d.", count, len(questions))
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
