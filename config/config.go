package config

import (
	//"fmt"
	"os"
)

type Settings struct {
	AuthKey			string
	ServiceUrl		string
	UserName		string
	Password		string
	UserToken		string
}

func NewConfiguration() * Settings {

	serviceUrl := getEnvironmentVariable("Article_Service_Url", "http://localhost:9000")
	authKey := getEnvironmentVariable("Ariticle_Server_AuthKey", "VIrPcAi4Rff0gBwdWklRl3ywMwgC6mZH") 

	configSettings := &Settings{
		authKey,
		serviceUrl,
		"",
		"",
		"",
	}

	return configSettings
}

func getEnvironmentVariable(envvar string, defaultValue string) string {

	variable := os.Getenv(envvar)

	if variable != "" {
		return variable
	}

	return defaultValue
}