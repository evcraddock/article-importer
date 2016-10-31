package tasks

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/evcraddock/article-importer/config"
	"github.com/evcraddock/article-importer/service"
)

type Task struct {
	service 			*service.HttpService
	articleLocation		string
}

func NewTask(settings *config.Settings) *Task {
	service := service.NewHttpService(settings.Auth)

	task := &Task{
		service,
		settings.ArticleLocation,
	}

	return task
}

func AskForStringValue(label string, defaultValue string, required bool) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		labelValue := label
		if defaultValue != "" {
			labelValue = label + " {" + defaultValue + "}"
		}

		if required {
			requiredtext := "*"
			labelValue = labelValue + " " + requiredtext
		}

		fmt.Printf("%s : ", labelValue)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		value := strings.Replace(response, "\n", "", -1)

		if len(value) == 0 {
			value = defaultValue
		}

		if required && strings.Trim(value, " ") == "" {
			fmt.Printf("")
			continue
		}

		return value
	}
}

func AskForCsv(label string, defaultValue []string) []string {
	csvstring := strings.Join(defaultValue, ", ")

	newcsv := AskForStringValue(label, csvstring, false)

	r := csv.NewReader(strings.NewReader(newcsv))
	stringArray, _ := r.Read()	
	return stringArray
}

func AskForDateValue(label string, defaultValue time.Time) time.Time {
	reader := bufio.NewReader(os.Stdin)
	dateValue := defaultValue

	for {
		fmt.Printf("%s {%d/%d/%d} : ", label, defaultValue.Month(), defaultValue.Day(), defaultValue.Year())

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		datestring := strings.Replace(response, "\n", "", -1)
		if len(datestring) == 0 {
			return defaultValue
		}

		dateValue, err = time.Parse("01/02/2006", datestring)
		if err != nil {
			fmt.Printf("Invalid Date, please try again {dd/mm/yyyy}\n")
			continue
		}

		return dateValue
	}
}

func getStringArray(value string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(value))
	return r.Read()
}