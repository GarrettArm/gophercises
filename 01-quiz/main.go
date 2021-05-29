package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type QandA struct {
	question string
	answer   string
}

type Score struct {
	corrects int
	totals   int
}

func parseFlags() (filepath *string, timeout *int) {
	filepath = flag.String("f", "problems.csv", "csv of q and a's")
	timeout = flag.Int("t", 30, "set a time limit for the quiz")
	flag.Parse()
	return filepath, timeout
}

func parseCSV(filepath string) []QandA {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	var records []QandA
	for {
		i, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		item := QandA{question: i[0], answer: i[1]}
		records = append(records, item)
	}
	return records
}

func doQuestion(qa QandA) bool {
	fmt.Printf("What is: %v?\n", qa.question)
	scanner := bufio.NewScanner(os.Stdin)
	var isCorrect bool
	for scanner.Scan() {
		if scanner.Text() == qa.answer {
			isCorrect = true
		} else {
			isCorrect = false
		}
		break
	}
	return isCorrect
}

func main() {
	score := Score{}
	filepath, timeout := parseFlags()
	QandAs := parseCSV(*filepath)
	timer := time.NewTimer(time.Duration(*timeout) * time.Second)

	answerChan := make(chan bool)
	for _, qa := range QandAs {
		go func() {
			answerChan <- doQuestion(qa)
		}()
		select {
		case correct := <-answerChan:
			if correct {
				score.corrects++
			}
			score.totals++
		case <-timer.C:
			fmt.Println("timer finished")
			fmt.Printf("%v out of %v correct\n", score.corrects, score.totals)
			return
		}
	}
}
