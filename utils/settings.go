package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var environments = map[string]string{
	"production":    "utils/prod.json",
	"preproduction": "utils/pre.json",
	"tests":         "../../utils/tests.json",
}

type Settings struct {
	PrivateKeyPath     string
	PublicKeyPath      string
	JWTExpirationDelta int
}

var settings Settings
var env = "preproduction"

func Init() {
	env = os.Getenv("GO_ENV")
	if env == "" {
		fmt.Println("Warning: Setting preproduction enviroment due to lack of GO_ENV value")
		env = "preproduction"
	}
	LoadSettingsByEnv(env)
}

func LoadSettingsByEnv(env string) Settings {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)
	}
	settings = Settings{}
	jsonErr := json.Unmarshal(content, &settings)
	if jsonErr != nil {
		fmt.Println("Error while parsing config file", jsonErr)
	}
	return settings
}

func GetEnvironment() string {
	return env
}

func GetSettings() Settings {
	if &settings == nil {
		Init()
	}
	return settings
}

func IsTestEnvironment() bool {
	return env == "tests"
}
