package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// apiStartURL is the starting fragment of API's URL.
var apiStartURL = "http://localhost:123/v1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI interface to the API quiz.",
	Long:  "CLI interface to the API quiz.",
}

// FetchQuestion fetches the question ID from API
func FetchQuestion(categoryId string) (io.ReadCloser, error) {
	if categoryIdnum, err := strconv.Atoi(categoryId); err != nil || categoryIdnum <= 0 {
		return nil, errors.New("question ID is not a valid ID")
	}

	questionURL := apiStartURL + "/questions/" + categoryId
	resp, err := http.Get(questionURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP status is not OK")
	}

	return resp.Body, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
