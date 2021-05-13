package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// startquizCmd represents the startquiz command
var startquizCmd = &cobra.Command{
	Use:   "startquiz",
	Short: "Fetch questions from server.",
	Long:  `Fetch question from API server. How to use it: cli.exe startquiz --category=2`,
	Run: func(cmd *cobra.Command, args []string) {
		categoryId, _ := cmd.Flags().GetString("category")

		questionReader, err := FetchQuestion(categoryId)
		if err != nil {
			os.Stderr.WriteString("Failed to fetch the question.")
			os.Exit(1)
		}
		defer questionReader.Close()

		respBody, err := ioutil.ReadAll(questionReader)
		if err != nil {
			os.Stderr.WriteString("Error reading response. " + err.Error())
			os.Exit(1)
		}

		fmt.Println("Response:\n", string(respBody))
	},
}

func init() {
	rootCmd.AddCommand(startquizCmd)

	startquizCmd.Flags().StringP("category", "c", "", "Category ID")
}
