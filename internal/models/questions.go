package models

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"youth-convention-2025/internal/constants"
)

type QA struct {
	Letter     string
	Difficulty string
	Question   string
	Answer     string
	Answers    []string
}

func (q *QA) PostProcess() {
	if q.Difficulty != constants.STR_HURDLE {
		return
	}

	for i, ans := range q.Answers {
		words := strings.SplitN(ans, " ", 2)
		q.Answers[i] = words[1]
	}
}

func (q *QA) Validate() {
	if q.Letter == "" {
		panic("No letter")
	}
	if q.Difficulty == "" {
		panic("No difficulty")
	}
	if q.Question == "" {
		panic("No question")
	}
	if len(q.Answers) == 0 && q.Answer == "" {
		panic("No answer")
	}
	if q.Answer == "" && len(q.Answers) == 0 {
		panic("No answers")
	}
}

func (q *QA) Print() {
	fmt.Println("Letter:", q.Letter)
	fmt.Println("Difficulty:", q.Difficulty)
	fmt.Println("Question:", q.Question)
	if q.Answer != "" {
		fmt.Println("Answer:", q.Answer)
	} else {
		fmt.Println("Answers:", q.Answers)
	}
	fmt.Println()
}

func QAsFromMarkdown(filepath string) []*QA {
	questions := make([]*QA, 0, 26*3)

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	const PREFIX_POUND = "# "
	const PREFIX_SPRINT = "1. Sprint"
	const PREFIX_MARATHON = "2. Marathon"
	const PREFIX_HURDLE = "3. Hurdle"
	const PREFIX_HURDLE_QUESTION = "3. Hurdle - "
	const PREFIX_QUESTION = "1. "
	const PREFIX_ANSWER = "1. Answer: "

	question := &QA{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if question.Letter == "" {
			if strings.HasPrefix(line, PREFIX_POUND) {
				left := strings.TrimPrefix(line, PREFIX_POUND)
				ascii := int(rune(left[0]))
				if ascii >= constants.ASCII_A && ascii <= constants.ASCII_Z {
					question.Letter = left
				}
			}
			continue
		}

		if question.Difficulty == "" {
			if strings.HasPrefix(line, PREFIX_SPRINT) {
				question.Difficulty = constants.STR_SPRINT
			}
			if strings.HasPrefix(line, PREFIX_MARATHON) {
				question.Difficulty = constants.STR_MARATHON
			}
			if strings.HasPrefix(line, PREFIX_HURDLE) {
				question.Difficulty = constants.STR_HURDLE
				question.Question = strings.TrimPrefix(line, PREFIX_HURDLE_QUESTION)
			}
			continue
		}

		if question.Question == "" {
			if strings.HasPrefix(line, PREFIX_QUESTION) {
				question.Question = strings.TrimPrefix(line, PREFIX_QUESTION)
			}
			continue
		}

		if question.Difficulty != constants.STR_HURDLE {
			if question.Answer == "" {
				if strings.HasPrefix(line, PREFIX_ANSWER) {
					question.Answer = strings.TrimPrefix(line, PREFIX_ANSWER)
				}
			}
		} else {
			question.Answers = append(question.Answers, line)
			for scanner.Scan() {
				line := scanner.Text()
				line = strings.TrimSpace(line)
				if len(line) == 0 {
					break
				}
				question.Answers = append(question.Answers, line)
			}
		}

		question.PostProcess()
		question.Validate()
		// question.Print()
		questions = append(questions, question)

		if question.Difficulty == constants.STR_HURDLE {
			question = &QA{}
		} else {
			question = &QA{
				Letter: question.Letter,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return questions
}

func GetQuestion(questions []*QA, letter string, difficulty string) *QA {
	var questionToUse *QA
	letterLow := strings.ToLower(letter)
	difficultyLow := strings.ToLower(difficulty)
	for _, question := range questions {
		if strings.ToLower(question.Letter) == letterLow && strings.ToLower(question.Difficulty) == difficultyLow {
			questionToUse = question
			break
		}
	}
	if questionToUse == nil {
		panic("No question found")
	}
	return questionToUse
}

func init() {
	_ = QAsFromMarkdown("./data/questions.md")
}
