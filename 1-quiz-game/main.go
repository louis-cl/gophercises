package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

type problem struct {
	question string
	answer   string
}

func createProblemList(data [][]string) []problem {
	var problems []problem
	for _, line := range data {
		problems = append(problems, problem{line[0], line[1]})
	}
	return problems
}

func play(problems []problem, timeLimitS int) {
	fmt.Printf("Answer each question in less than %d seconds.\nPress Enter to start !", timeLimitS)
	fmt.Scanln()

	scanner := bufio.NewScanner(os.Stdin)

	correctAnswers := 0
	// timer_duration := time.Duration(timeLimitS) * time.Second
	// timer := time.NewTimer(timer_duration)

	answerCh := make(chan string)

	for i, prob := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, prob.question)
		go func() {
			scanner.Scan()
			answerCh <- scanner.Text()
		}()
		select {
		case answer := <-answerCh:
			if answer == prob.answer {
				correctAnswers++
			}
		}
		// timer.Reset(timer_duration)
	}
	fmt.Printf("You scored %d out of %d", correctAnswers, len(problems))
}

func main() {
	var (
		csvFile   = flag.String("csv", "problems.csv", "a csv file in the format of question,answer")
		timeLimit = flag.Int("limit", 30, "the time limit for the quiz in seconds")
	)
	flag.Parse()
	// read file
	f, err := os.Open(*csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// parse csv
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	// convert data to problems
	problems := createProblemList(data)
	play(problems, *timeLimit)
}
