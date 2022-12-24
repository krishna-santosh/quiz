package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func makeCSV() {
	fmt.Println("CSV Generator")
	fmt.Print("What do you want your CSV file to be called: ")

	fileName := getInput()

	if !strings.HasSuffix(fileName, ".csv") {
		fileName += ".csv"
	}

	file, err := os.Create(strings.TrimSpace(fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	fmt.Print("Enter the number of questions: ")
	qcount, err := strconv.Atoi(getInput())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Would you like to hash your answers? [Y/n]: ")
	hashResp := getInput()
	toHash := true
	if hashResp == "N" || hashResp == "n" {
		toHash = false
	}

	for i := 1; i <= qcount; i++ {
		r := questin_answer_input(i, toHash)
		err := writer.Write(r)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("\nSuccessfully created %s, and added %d questions to it.\n", fileName, qcount)
	if toHash {
		fmt.Printf("\nTo play with %s run `quiz -csv %s -hash`\n", fileName, fileName)
	} else {
		fmt.Printf("\nTo play with %s run `quiz -csv %s`\n", fileName, fileName)
	}
}

func questin_answer_input(question_no int, hash bool) []string {
	fmt.Printf("Question %d: ", question_no)
	qinput := getInput()

	fmt.Print("Answer : ")
	ainput := getInput()

	if hash {
		ainput = fmt.Sprintf("%x", sha256.Sum256([]byte(ainput)))
	}

	return []string{qinput, ainput}
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return strings.TrimSpace(input)
}
