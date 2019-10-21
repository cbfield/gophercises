package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func checkErr(err error, msg string) {
	if err != nil {
		fmt.Printf(msg)
		os.Exit(1)
	}
}

type problem struct {
	question string
	answer   string
}

func parseProblems(lines [][]string) []problem {
	returnValue := make([]problem, len(lines))
	for i, line := range lines {
		returnValue[i] = problem{
			question: strings.TrimSpace(line[0]),
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return returnValue
}

func serveProblem(i int, problem *problem) bool {
	fmt.Printf("Problem #%d: %s = ", i+1, problem.question)
	var answer string
	fmt.Scanf("%s\n", &answer)
	return answer == problem.answer
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "csv question,answer pairs")
	timeLimit := flag.Int("limit", 15, "quiz time limit (seconds)")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	checkErr(err, fmt.Sprintf("Failed to open file: %s\n", *csvFileName))

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	checkErr(err, fmt.Sprintf("Failed to parse file: %s", *csvFileName))

	problems := parseProblems(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, problem := range problems {
		select {
		case <-timer.C:
			fmt.Printf("Score: %d/%d", correct, len(problems))
			return
		default:
			outcome := serveProblem(i, &problem)
			if outcome {
				correct++
			}
		}
	}
}
