package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
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
	fmt.Printf("Answer all questions in less than %d seconds.\nPress Enter to start !", timeLimitS)
	fmt.Scanln()

	scanner := bufio.NewScanner(os.Stdin)

	correctAnswers := 0

	timer := time.NewTimer(time.Duration(timeLimitS) * time.Second)
	answerCh := make(chan string)

out:
	for i, prob := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, prob.question)
		go func() {
			scanner.Scan()
			answerCh <- strings.ToLower(strings.TrimSpace(scanner.Text()))
		}()
		select {
		case <-timer.C:
			fmt.Println("\nYou ran out of time !")
			break out
		case answer := <-answerCh:
			if answer == prob.answer {
				correctAnswers++
			}
		}
	}
	fmt.Printf("You scored %d out of %d", correctAnswers, len(problems))
}

func main() {
	var (
		csvFile   = flag.String("csv", "problems.csv", "a csv file in the format of question,answer")
		timeLimit = flag.Int("limit", 30, "the time limit for the quiz in seconds")
		shuffle   = flag.Bool("shuffle", false, "set to shuffle questions randomly")
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
	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	}
	play(problems, *timeLimit)
}
