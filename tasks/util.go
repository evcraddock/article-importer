package tasks

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func AskForStringValue(label string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		labelValue := label
		if defaultValue != "" {
			labelValue = label + " <" + defaultValue + ">"
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

		return value
	}
}

func AskForDateValue(label string) time.Time {
	reader := bufio.NewReader(os.Stdin)
	now := time.Now()
	dateValue := time.Now()

	for {
		fmt.Printf("%s (%d/%d/%d) : ", label, now.Month(), now.Day(), now.Year())

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		datestring := strings.Replace(response, "\n", "", -1)
		if len(datestring) == 0 {
			return now
		}

		dateValue, err = time.Parse("01/02/2006", datestring)
		if err != nil {
			fmt.Printf("Invalid Date, please try again {dd/mm/yyyy}\n")
			continue
		}

		return dateValue
	}
}