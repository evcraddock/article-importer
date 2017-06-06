package config

import (
	"os"
)

type Settings struct {
	Auth            Authorization
	ArticleLocation string
}

type Authorization struct {
	AuthKey    string
	ServiceUrl string
	UserName   string
	Password   string
}

func NewConfiguration() *Settings {

	serviceUrl := getEnvironmentVariable("Article_Service_Url", "http://localhost:9000")
	authKey := getEnvironmentVariable("Ariticle_Server_AuthKey", "VIrPcAi4Rff0gBwdWklRl3ywMwgC6mZH")
	articleLocation := getEnvironmentVariable("Article_Location", "articles/")

	authSettings := &Authorization{
		authKey,
		serviceUrl,
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
