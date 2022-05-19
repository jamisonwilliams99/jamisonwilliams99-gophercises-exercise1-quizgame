package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func OpenCSV() (*os.File, error) {
	var f *os.File
	var err error

	if len(os.Args) < 2 {
		f, err = os.Open("problems.csv")

		if err != nil {
			return nil, fmt.Errorf("Unable to open file '%v'\n", os.Args[3])
		}
	} else {
		if os.Args[1] == "-r" {
			if len(os.Args) < 3 {
				return nil, fmt.Errorf("Missing file name argument")
			}

			if os.Args[2][len(os.Args[2])-4:] != ".csv" {
				return nil, fmt.Errorf("Specified file must be a .csv")
			}

			f, err = os.Open(os.Args[2])

			if err != nil {
				return nil, fmt.Errorf("Unable to open file '%v'\n", os.Args[3])
			}
		} else {
			return nil, fmt.Errorf("Unknown argument '%v'\n", os.Args[1])
		}
	}

	return f, nil
}

func promptUser(operation string) string {
	var userAnswer string

	fmt.Printf("%v: ", operation)
	fmt.Scanln(&userAnswer)
	fmt.Println()

	return userAnswer
}

func quiz(f *os.File) (int, int) {
	csvReader := csv.NewReader(f)
	questions, _ := csvReader.ReadAll()

	correctAnswers := 0
	totalQuestions := len(questions) - 1

	timer := time.NewTimer(30 * time.Second)

	finished := make(chan bool)

	go func() {
		for _, question := range questions {

			operation, answer := question[0], question[1]

			userAnswer := promptUser(operation)

			if userAnswer == answer {
				correctAnswers++
			}
		}
		finished <- true
	}()

	// selects whichever channel finishes first
	select {
	case <-finished:
	case <-timer.C:
		fmt.Println("time is up!")
	}

	return correctAnswers, totalQuestions
}

func main() {
	f, err := OpenCSV()

	if err != nil {
		log.Fatal(err)
	}

	correctAnswers, totalQuestions := quiz(f)
	score := (float64(correctAnswers) / float64(totalQuestions)) * 100

	fmt.Printf("You got %v/%v correct! (score: %.2f)\n", correctAnswers, totalQuestions, score)

}
