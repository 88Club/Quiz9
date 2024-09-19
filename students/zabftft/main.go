package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	ErrIncorrectFileFormat = errors.New("invalid file format")
	ErrTimeLimitIsNegative = errors.New("negative time limit")
)

const (
	FileFlag         = "file"
	FileName         = "problems.csv"
	TimeLimitFlag    = "limit"
	TimeLimitSeconds = 3
	ShuffleFlag      = "shuffle"
)

type gameResult struct {
	total   int
	correct int
	err     error
}

func main() {
	fileName := flag.String(FileFlag, FileName, "filename where questions and asnwers are placed for the quiz")
	timeLimit := flag.Int(TimeLimitFlag, TimeLimitSeconds, "time limitation for the quiz in seconds")
	shuffle := flag.Bool(ShuffleFlag, false, "shuffle quiz questions")
	flag.Parse()

	if err := validateFlags(*fileName, *timeLimit, *shuffle); err != nil {
		log.Fatal(err.Error())
	}

	fileContent, err := getFileContent(*fileName)

	if err != nil {
		log.Fatal(err.Error())
	}

	if *shuffle == true {
		fileContent = shuffleQuestions(fileContent)
	}

	quizDuration := time.Duration(*timeLimit) * time.Second
	runQuiz(fileContent, quizDuration)
}

func validateFlags(fileName string, timeLimit int, shuffle bool) error {
	if !strings.Contains(fileName, ".csv") {
		return ErrIncorrectFileFormat
	}

	if timeLimit <= 0 {
		return ErrTimeLimitIsNegative
	}

	return nil
}

func getFileContent(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	return r.ReadAll()
}

func shuffleQuestions(fileContent [][]string) [][]string {
	rand.Shuffle(len(fileContent), func(i, j int) {
		fileContent[i], fileContent[j] = fileContent[j], fileContent[i]
	})
	return fileContent
}

func runQuiz(data [][]string, timeLimit time.Duration) error {
	r := bufio.NewReader(os.Stdin)
	gc := make(chan gameResult)

	entryPrompt(r)

	timer := time.NewTimer(timeLimit)
	go askQuestions(r, data, gc)

	var result gameResult
	for {
		select {
		case <-timer.C:
			close(gc)
			fmt.Printf("\nTotal Questions: %d. Correct Answers: %d", result.total, result.correct)
			return nil
		case result = <-gc:
			if result.err != nil {
				return result.err
			}
		}
	}
}

func entryPrompt(r *bufio.Reader) error {
	fmt.Print("Welcome to Quiz!\n\n")
	fmt.Print("Press Enter when you are ready\n")
	_, err := r.ReadString('\n')

	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}

func askQuestions(r *bufio.Reader, data [][]string, gc chan<- gameResult) {
	var totalQuestions int
	var correctAnswers int

	for i, question := range data {
		quizQuestion, correctAnswer := question[0], question[1]
		fmt.Printf("Question #%d: %s\n", i+1, quizQuestion)
		fmt.Print("Your answer: ")

		answer, err := r.ReadString('\n')

		if err != nil {
			gc <- gameResult{err: err}
			break
		}

		answer = cleanAnswer(answer)

		if answer == correctAnswer {
			correctAnswers++
		}
		totalQuestions++

		gc <- gameResult{total: totalQuestions, correct: correctAnswers}
	}
}

func cleanAnswer(answer string) string {
	return strings.TrimSpace(strings.ToLower(answer))
}
