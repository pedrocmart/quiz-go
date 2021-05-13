# quiz-go

# A quiz API server and CLI client in Golang
The task is to build a super simple quiz with a few questions and a few alternatives for each question. With one correct answer.

# Database
Since no database is required, I tried to reproduce a database with 4 categories sets:<br/>
1- Movies<br/>
2- Geography<br/>
3- History<br/>
4- General Knowledge

These IDs (1, 2, 3 or 4) will be required to run the CLI, depending on which chosen category you'll answer the quiz.

## To run the API

CD to ./src/api/ and exec: go run main.go errorHandler.go questionHandler.go data.go

## To run the CLI
CD to ./src/cli/ and exec: go install . And after that: ./bin/cli.exe start --category=2<br/>
or just ./bin/cli.exe start -c=2


## User stories/Use cases: <br/>
-User should be able to get questions with a number of answers<br/>
-User should be able to select just one answer per question.<br/>
-User should be able to answer all the questions and then post his/hers answers and get back how many correct answers they had, displayed to the user.<br/>
-User should see how good he/she did compare to others that have taken the quiz, "You were better than 60% of all quizzers"<br/>
