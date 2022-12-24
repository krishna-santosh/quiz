package main

import (
	"crypto/sha256"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {

	csvFileName := flag.String("csv", "problems.csv", "a csv file in 'question,answer' format.")
	timeLimit := flag.Int("time", 10, "time limit in seconds")
	hashed := flag.Bool("hash", false, "to work with sha256 hashed answers")
	makecsv := flag.Bool("make-csv", false, "to create a new csv file for quiz")
	flag.Parse()

	if flag.NFlag() == 1 {
		if *makecsv {
			makeCSV()
			os.Exit(0)
		}
	}

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open %s", *csvFileName))
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	fmt.Println("Welcome to Quiz game.")

	var correct_answer_count = 0
	answer_channel := make(chan string)

	fmt.Printf("\nYou have %d seconds to finish the quiz.\n", *timeLimit)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for index, problem := range parseLines(lines) {
		fmt.Printf("\nQuestion #%d/%d\n%s = ", index+1, len(parseLines((lines))), problem.question)

		go func() {
			var answer_input string
			fmt.Scanln(&answer_input)
			if *hashed {
				hashed_input_ans := fmt.Sprintf("%x", sha256.Sum256([]byte(answer_input)))
				answer_channel <- hashed_input_ans
			} else {
				answer_channel <- answer_input
			}
		}()

		select {
		case <-timer.C:
			fmt.Printf("\n\nTime out")
			fmt.Printf("\nYou scored %d out of %d\n", correct_answer_count, len(parseLines((lines))))
			return

		case answer := <-answer_channel:
			if answer == problem.answer {
				correct_answer_count++
			}
		}
	}

	fmt.Printf("\nYou scored %d out of %d\n", correct_answer_count, len(parseLines((lines))))
}

func parseLines(lines [][]string) []problem {
	array_of_problems := make([]problem, len(lines))

	for index, line := range lines {
		array_of_problems[index] = problem{
			question: strings.TrimSpace(line[0]),
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return array_of_problems
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
