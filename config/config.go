package config

import (
	"github.com/tkanos/gonfig"
	"os"
)

type Configuration struct {
	DB_USERNAME       string
	DB_PASSWORD       string
	DB_PORT           string
	DB_HOST           string
	DB_NAME           string
	APP_SECRET        string
	MAILGUN_DOMAIN    string
	MAILGUN_KEY       string
	STREAMER_SERVER   string
	SMS_API_URL       string
	SMS_API_USER_NAME string
	SMS_API_PASSWORD  string
	SMS_API_SOURCE    string
	SAAS_BASE         bool
}

func config() Configuration {
	var configFileName string
	configuration := Configuration{}

	if os.Getenv("GO_ENV") == "production" {
		configFileName = "prod.json"
	} else
	{
		configFileName = "prod.json"
	}
	gonfig.GetConf("config/"+configFileName, &configuration)

	return configuration
}

func GetHmacKey() string {
	config := config()
	return config.APP_SECRET
}
func GetHmacSignKey() []byte {
	config := config()
	Secret := []byte(config.APP_SECRET)
	return Secret
}
