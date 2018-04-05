package tasks

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/evcraddock/article-importer/config"
	"github.com/evcraddock/article-importer/service"
	"golang.org/x/crypto/ssh/terminal"
)

//Task stores task information
type Task struct {
	service *service.HTTPService
}

//NewTask creates new instance of a Task
func NewTask(settings *config.Settings) *Task {
	service := service.NewHTTPService(settings.Auth)

	task := &Task{
		service,
	}

	return task
}

//AskForStringValue prompts user for a string value
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

//AskForHiddenStringValue prompts the user for a value which should not be displayed on the screen
func AskForHiddenStringValue(label string, defaultValue string, required bool) string {
	for {
		labelValue := label
		if defaultValue != "" {
			labelValue = label + " { ******** }"
		}

		if required {
			requiredtext := "*"
			labelValue = labelValue + " " + requiredtext
		}

		fmt.Printf("%s : ", labelValue)
		byteHidden, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Printf("\n")

		if err != nil {
			log.Fatal(err)
		}

		hiddentext := string(byteHidden)
		return hiddentext
	}
}

//AskForCSV prompts user for value seperated by commas
func AskForCSV(label string, defaultValue []string) []string {
	csvstring := removeWhiteSpace(strings.Join(defaultValue, ", "))

	newcsv := AskForStringValue(label, csvstring, false)

	r := csv.NewReader(strings.NewReader(newcsv))
	stringArray, _ := r.Read()
	return stringArray
}

//AskForDateValue prompts user for a date
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

//GetFileName returns the filename from a path
func GetFileName(value, delimiter string) string {
	fullarray := strings.Split(value, delimiter)
	lastvalue := fullarray[len(fullarray)-1]

	if lastvalue == "" {
		return fullarray[len(fullarray)-2]
	}

	return lastvalue
}

func getStringArray(value string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(value))
	return r.Read()
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

func removeWhiteSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
