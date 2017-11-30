package config

import (
	"os"
)

//Settings object for storing settings
type Settings struct {
	Auth            Authorization
	ArticleLocation string
}

//Authorization object for keeping credentials
type Authorization struct {
	AuthKey    string
	ServiceURL string
	UserName   string
	Password   string
}

//NewConfiguration creates a new Settings instance
func NewConfiguration() *Settings {
	serviceURL := getEnvironmentVariable("Article_Service_Url", "http://localhost:9000")
	authKey := getEnvironmentVariable("Ariticle_Server_AuthKey", "VIrPcAi4Rff0gBwdWklRl3ywMwgC6mZH")
	articleLocation := getEnvironmentVariable("Article_Location", "~/articles/")

	authSettings := &Authorization{
		authKey,
		serviceURL,
		"",
		"",
	}

	configSettings := &Settings{
		*authSettings,
		articleLocation,
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
