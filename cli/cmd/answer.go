package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const (
	c1 = "MOVIES"
	c2 = "GEOGRAPHY"
	c3 = "HISTORY"
	c4 = "GENERAL KNOWLEDGE"
)

// Question's structure
type Question struct {
	Question string    `json:"question"`
	Answers  [4]string `json:"answers"`
}

// Answer's structure
type Answer struct {
	AnswerID int `json:"answerId"`
}

type Answers [5]Answer

// startCmd represents the answer command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Quiz",
	Long: `Start the Quiz and take questions from API server.

	How to use it: cli.exe start --category=1 , where --category is an ID of the category.
	Short usage: cli.exe start -c=1 .
	[1] Movies [2] Geography [3] History [4] General Knowledge`,

	Run: func(cmd *cobra.Command, args []string) {

		// read category ID from the argument
		categoryId, _ := cmd.Flags().GetString("category")

		// fetch the question from API
		fmt.Print("Fetching questions from API...\n\n")
		questionReader, err := FetchQuestion(categoryId)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		defer questionReader.Close()

		var questions [5]Question
		err = json.NewDecoder(questionReader).Decode(&questions)
		if err != nil {
			os.Stderr.WriteString("Failed to parse question's JSON")
			os.Exit(1)
		}

		var answerW string
		var answers [len(questions)]Answer
		var questionStr strings.Builder

		printHeader(categoryId)

		// Loop through the questions and output to the user with options
		for i, question := range questions {
			answerW = ""
			questionStr.WriteString(question.Question + ":\n")
			answersLen := len(question.Answers)
			for j, answer := range question.Answers {
				questionStr.WriteString("[")
				questionStr.WriteString(strconv.Itoa(j + 1))
				questionStr.WriteString("] ")
				questionStr.WriteString(answer)
				questionStr.WriteString("\n")
			}
			questionStr.WriteString("Enter the right answer's number: ")
			fmt.Println(questionStr.String())
			questionStr.Reset()

			// read user input
			fmt.Scanln(&answerW)
			fmt.Print("\n")

			// validate user input
			answerNum, err := strconv.Atoi(answerW)
			if err != nil {
				os.Stderr.WriteString("Answer must be a number")
				os.Exit(1)
			}
			if answerNum <= 0 || answerNum > answersLen {
				answersLenStr := strconv.Itoa(answersLen)
				os.Stderr.WriteString("Answer must be in the range 1-" + answersLenStr)
				os.Exit(1)
			}

			answers[i] = Answer{
				AnswerID: answerNum - 1, // decrementing to match the same database index
			}
		}

		fmt.Print("Sending your answers to server...\n\n")

		// sending our answers and outputting results to the CLI
		answersResultReader, err := sendAnswer(categoryId, answers)
		if err != nil {
			os.Stderr.WriteString("Failed to send answers. " + err.Error())
			os.Exit(1)
		}
		defer answersResultReader.Close()

		respBody, err := ioutil.ReadAll(answersResultReader)
		if err != nil {
			os.Stderr.WriteString("Error reading response. " + err.Error())
			os.Exit(1)
		}

		printResults(respBody)
	},
}

// Send answers to the API
func sendAnswer(categoryId string, answers Answers) (io.ReadCloser, error) {
	answersJSON, err := json.Marshal(answers)
	if err != nil {
		os.Stderr.WriteString("Failed to JSON encode our answers. " + err.Error())
		os.Exit(1)
	}

	answerURL := apiStartURL + "/questions/" + categoryId + "/answers"
	resp, err := http.Post(answerURL, "application/json", bytes.NewBuffer(answersJSON))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP status is not OK." + answerURL)
	}

	return resp.Body, nil
}

func printHeader(categoryId string) {
	fmt.Print("-------------- SUPER SIMPLE QUIZ! --------------\n\n")
	fmt.Print("You chose the category: " + getCategory(categoryId))
	fmt.Print("\n\nGood Luck! \n\n\n")
}

func printResults(respBody []byte) {
	var result map[string]interface{}
	json.Unmarshal([]byte(respBody), &result)

	fmt.Println("-------------- Your results: --------------")
	fmt.Println("- Correct Answers: ", result["correctAnswers"])
	fmt.Println("- Total users who answered: ", result["answeredUsers"])
	fmt.Print("- You were better than ", result["comparingToOthers"])
	fmt.Print("% of all quizzers \n")
	fmt.Println("---------------------------------------------")
	fmt.Print("\n\n")
}

func getCategory(categoryId string) string {
	switch categoryId {
	case "1":
		return c1
	case "2":
		return c2
	case "3":
		return c3
	default:
		return c4
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("category", "c", "", "Category ID | [1] Movies [2] Geography [3] History [4] General Knowledge")
}
