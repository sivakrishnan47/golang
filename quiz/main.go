package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "csv file of format question and answer")
	timer := flag.Int("timer", 10, "add time in seconds, default time set is 30 seconds")
	flag.Parse()
	quizTimer := time.NewTimer(time.Duration(*timer) * time.Second)
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("File cannot be opened %s", *csvFilename))
	}
	csvReader := csv.NewReader(file)
	lines, err := csvReader.ReadAll()
	if err != nil {
		exit("failed to to read csv file")
	}

	problems := parseLines(lines)
	if len(problems) > 0 {

		correct := 0
		answerCh := make(chan string)
	problemloop:
		for index, problem := range problems {
			fmt.Printf("%d.%s =", index+1, problem.question)
			go func() {
				var answer string
				fmt.Scanf("%s\n", &answer)
				answerCh <- answer
			}()

			select {
			case <-quizTimer.C:
				fmt.Println("\ntest timed out")
				break problemloop
			case answer := <-answerCh:
				if answer == problem.answer {
					correct++
				}
			}
		}
		fmt.Printf("%d questions answered correctly of %d\n", correct, len(problems))
	}

}

func parseLines(lines [][]string) (problems []problem) {
	problems = make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem{
			question: line[0],
			answer:   line[1],
		}
	}
	return
}

func exit(errMsg string) {
	fmt.Println(errMsg)
	os.Exit(1)
}
