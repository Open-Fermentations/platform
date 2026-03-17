package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type AppEnv string

var AppEnvEnum = &struct {
	Dev  AppEnv
	Prod AppEnv
}{
	Dev:  "dev",
	Prod: "prod",
}

type DatabaseEnv struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	Schema   string
}

type Env struct {
	Port     int
	AppEnv   AppEnv
	Database DatabaseEnv
}

var env *Env

func GetEnv() *Env {
	if env == nil {
		RefreshEnvironmentVariables()
	}

	return env
}

func RefreshEnvironmentVariables() {
	env = &Env{
		Port:   getIntValue("PORT"),
		AppEnv: handleAppEnv("APP_ENV"),
		Database: DatabaseEnv{
			Host:     getStringValue("DB_HOST"),
			Port:     getStringValue("DB_PORT"),
			User:     getStringValue("DB_USERNAME"),
			Password: getStringValue("DB_PASSWORD"),
			DbName:   getStringValue("DB_DATABASE"),
			Schema:   getStringValue("DB_SCHEMA"),
		},
	}
}

func getStringValue(key string) string {
	e := os.Getenv(key)
	if e == "" {
		panic(fmt.Sprintf("%v was not defined in environment variables", key))
	}

	return e
}

func getIntValue(key string) int {
	e := os.Getenv(key)
	value, err := strconv.Atoi(e)
	if err != nil {
		panic(err)
	}

	return value
}

func handleAppEnv(key string) AppEnv {
	switch key {
	case string(AppEnvEnum.Dev):
		return AppEnvEnum.Dev
	case string(AppEnvEnum.Prod):
		return AppEnvEnum.Prod
	}

	log.Printf("No environment variable was set in %v defaulting to %v", key, AppEnvEnum.Dev)
	return AppEnvEnum.Dev
}
