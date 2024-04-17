package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

// Fisher-Yates shuffle algorithm
func Shuffle(data []problem) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(data); i++ {
		r := random.Intn(i + 1)
		data[i], data[r] = data[r], data[i]
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseLines(lines [][]string) []problem {
	res := make([]problem, len(lines))

	for i, line := range lines {
		res[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return res
}

func main() {
	csvFileName := flag.String(
		"csv",
		"problems.csv",
		"a csv file in the format of 'question, answer'",
	)
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")

	shuffle := flag.Bool("shuffle", false, "shuffle order of the questions")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFileName))
	}

	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the CSV file.")
	}
	problems := parseLines(lines)

	if *shuffle {
		Shuffle(problems)
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
problemLoop:
	for index, problem := range problems {

		fmt.Printf("Problem #%d: %s = ", index+1, problem.q)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answer = strings.TrimSpace(answer)
			answer = strings.ToUpper(answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Println()
			break problemLoop
		case answer := <-answerCh:

			if problem.a == answer {
				correct++
			}

		}

	}
	// <-timer.C

	fmt.Printf("\nYou scored %d out of %d", correct, len(problems))
}
