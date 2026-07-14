package env

import (
	"log/slog"
	"open-fermentations/internal/logging"
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

type LogLevel string

var LogLevelEnum = &struct {
	Debug   LogLevel
	Info    LogLevel
	Warning LogLevel
	Error   LogLevel
	None    LogLevel
}{
	Debug:   "debug",
	Info:    "info",
	Warning: "warn",
	Error:   "error",
	None:    "none",
}

type DatabaseEnv struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	Schema   string
}

type MqttEnv struct {
	Host     string
	Port     string
	WsPort   string
	User     string
	Password string
	ClientID string
}

type JwtEnv struct {
	Key    string
	Issuer string
}

type Env struct {
	Port         int
	AppEnv       AppEnv
	LogLevel     LogLevel
	Database     DatabaseEnv
	Mqtt         MqttEnv
	Jwt          JwtEnv
	CookieSecure bool
}

var env *Env

func GetEnv() *Env {
	if env == nil {
		env = &Env{}
		RefreshEnvironmentVariables()
	}

	return env
}

func RefreshEnvironmentVariables() {
	env.Port = getIntValue("PORT", 8080)
	env.AppEnv = handleAppEnv("APP_ENV")
	env.LogLevel = handleLogLevel("LOG_LEVEL")

	env.Jwt.Key = getStringValue("JWT_KEY", "")
	env.Jwt.Issuer = getStringValue("JWT_ISSUER", "open-fermentations")

	env.CookieSecure = getBoolValue("COOKIE_SECURE", true)

	env.Database.Host = getStringValue("DB_HOST", "localhost")
	env.Database.Port = getStringValue("DB_PORT", "5432")
	env.Database.User = getStringValue("DB_USERNAME", "ferment")
	env.Database.Password = getStringValue("DB_PASSWORD", "password1234")
	env.Database.DbName = getStringValue("DB_DATABASE", "open-fermentations")
	env.Database.Schema = getStringValue("DB_SCHEMA", "public")

	env.Mqtt.Port = getStringValue("MQTT_PORT", "todo")
	env.Mqtt.WsPort = getStringValue("MQTT_WS_PORT", "todo")
	env.Mqtt.Host = getStringValue("MQTT_HOST", "mqtt")
	env.Mqtt.User = getStringValue("MQTT_USER", "platform")
	env.Mqtt.Password = getStringValue("MQTT_PASSWORD", "password1234")
	env.Mqtt.ClientID = getStringValue("MQTT_CLIENT_ID", "open-fermentations")
}

func getStringValue(key string, def string) string {
	e, ok := os.LookupEnv(key)
	if !ok {
		slog.Warn("environment variable default",
			slog.String("key", key),
			slog.String("default", def))
		return def
	}

	return e
}

func getIntValue(key string, def int) int {
	e, ok := os.LookupEnv(key)
	if !ok {
		slog.Warn("environment variable default",
			slog.String("key", key),
			slog.Int("default", def))
		return def
	}

	value, err := strconv.Atoi(e)
	if err != nil {
		slog.Warn("environment variable: parse int",
			slog.String("key", key),
			slog.String("value", e),
			slog.Int("default", def),
			logging.Err(err))
		return def
	}

	return value
}

func getBoolValue(key string, def bool) bool {
	e, ok := os.LookupEnv(key)
	if !ok {
		slog.Warn("environment variable default",
			slog.String("key", key),
			slog.Bool("default", def))
		return def
	}

	val, err := strconv.ParseBool(e)
	if err != nil {
		slog.Warn("environment variable: parse bool",
			slog.String("key", key),
			slog.String("value", e),
			slog.Bool("default", def),
			logging.Err(err))
		return def
	}

	return val
}

func handleAppEnv(key string) AppEnv {
	def := AppEnvEnum.Dev
	value := os.Getenv(key)
	switch value {
	case string(AppEnvEnum.Dev):
		return AppEnvEnum.Dev
	case string(AppEnvEnum.Prod):
		return AppEnvEnum.Prod
	}

	slog.Info("Defaulting",
		slog.String("key", key),
		slog.String("default", string(def)))
	return def
}

func handleLogLevel(key string) LogLevel {
	value := os.Getenv(key)
	switch value {
	case string(LogLevelEnum.Debug):
		return LogLevelEnum.Debug
	case string(LogLevelEnum.Info):
		return LogLevelEnum.Info
	case string(LogLevelEnum.Warning):
		return LogLevelEnum.Warning
	case string(LogLevelEnum.Error):
		return LogLevelEnum.Error
	default:
		slog.Info("Defaulting",
			slog.String("key", key),
			slog.String("default", string(LogLevelEnum.Info)))
		return LogLevelEnum.Info
	}
}
