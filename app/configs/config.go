package configs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

const (
	prod = "production"
	dev  = "development"
	stg  = "staging"
)

// Config object
type Config struct {
	Env                       string      `env:"ENV"`
	Mysql                     MysqlConfig `json:"mysql"`
	JWTAccessSecret           string      `env:"JWT_ACCESS_SIGN_KEY"`
	JWTRefreshSecret          string      `env:"JWT_REFRESH_SIGN_KEY"`
	JWTIssuer                 string      `env:"JWT_ISSUER"`
	Host                      string      `env:"APP_HOST"`
	Port                      string      `env:"APP_PORT"`
	DbHost                    string      `env:"DB_HOST"`
	DbPort                    string      `env:"DB_PORT"`
	DbDriver                  string      `env:"DB_DRIVER"`
	DbUser                    string      `env:"DB_USER"`
	DbPassword                string      `env:"DB_PASSWORD"`
	DbName                    string      `env:"DB_NAME"`
	GoogleOauthClientId       string      `env:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOauthClientSecret   string      `env:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOauthRedirectionUrl string      `env:"GOOGLE_OAUTH_REDIRECTION_URL"`

	/*RedisHost                   string      `env:"REDIS_HOST"`
	RedisPost                   string      `env:"REDIS_PORT"`
	RedisDB                     string      `env:"REDIS_DB"`
	RedisPassword               string      `env:"REDIS_PASSWORD"`*/
}

// IsProd Checks if env is production
func (c Config) IsProd() bool {
	return c.Env == prod
}

func (c Config) IsDev() bool {
	return c.Env == dev
}

func (c Config) IsStg() bool {
	return c.Env == stg
}

// LoadLocalConfig gets config from .env
func LoadLocalConfig() {
	requiredEnvVars := getRequiredEnvVars(Config{})

	var missingVars []string

	currentPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	environmentPath := filepath.Join(currentPath, ".env")

	if err := godotenv.Load(environmentPath); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	for _, envVar := range requiredEnvVars {
		if getEnv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		log.Fatalf("missing environment variables: %s", strings.Join(missingVars, ", "))
	}
}

func LoadSSMConfig() {
	requiredEnvVars := getRequiredEnvVars(Config{})

	var missingVars []string

	// Create an AWS session
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	ssmClient := ssm.NewFromConfig(cfg)
	envPath := os.Getenv("ENV_LOAD_PATH")
	if len(envPath) == 0 {
		log.Fatalln("ENV_LOAD_PATH env missing 🙇")
	}

	for _, envVar := range requiredEnvVars {
		isDecryption := true
		// Use AWS SSM to get the parameter value
		if os.Getenv(envVar) == "" {
			paramName := fmt.Sprintf("%s/%s", envPath, envVar)

			result, err := ssmClient.GetParameter(context.TODO(), &ssm.GetParameterInput{
				Name:           &paramName,
				WithDecryption: &isDecryption,
			})

			if err != nil {
				log.Fatalf("Error retrieving parameter %s: %s", envVar, err)
			}

			paramValue := *result.Parameter.Value

			fmt.Printf("Fetching param name : %s value :%v\n", paramName, paramValue)
			// Set the environment variable
			err = os.Setenv(envVar, paramValue)
			if err != nil {
				log.Fatalf("error while setting the error: %v", err)
			}

			if paramValue == "" {
				missingVars = append(missingVars, envVar)
			}
		}
	}

	if len(missingVars) > 0 {
		log.Fatalf("missing environment variables: %s", strings.Join(missingVars, ", "))
	}
}

func getRequiredEnvVars(cfgStruct interface{}) []string {
	var envVars []string
	val := reflect.ValueOf(cfgStruct)

	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("env")
		if tag != "" {
			envVars = append(envVars, tag)
		}
	}

	return envVars
}

// GetConfig gets all config for the application
func GetConfig() Config {
	return Config{
		Env:   getEnv("ENV"),
		Mysql: GetMysqlConfig(),
		// Mailgun:   GetMailgunConfig(),
		JWTAccessSecret:           getEnv("JWT_ACCESS_SIGN_KEY"),
		JWTRefreshSecret:          getEnv("JWT_REFRESH_SIGN_KEY"),
		JWTIssuer:                 getEnv("JWT_ISSUER"),
		Host:                      getEnv("APP_HOST"),
		Port:                      getEnv("APP_PORT"),
		DbHost:                    getEnv("DB_HOST"),
		DbPort:                    getEnv("DB_PORT"),
		DbDriver:                  getEnv("DB_DRIVER"),
		DbUser:                    getEnv("DB_USER"),
		DbPassword:                getEnv("DB_PASSWORD"),
		DbName:                    getEnv("DB_NAME"),
		GoogleOauthClientId:       getEnv("GOOGLE_OAUTH_CLIENT_ID"),
		GoogleOauthClientSecret:   getEnv("GOOGLE_OAUTH_CLIENT_SECRET"),
		GoogleOauthRedirectionUrl: getEnv("GOOGLE_OAUTH_REDIRECTION_URL"),
		/*RedisHost:                   getEnv("REDIS_HOST"),
		RedisPost:                   getEnv("REDIS_PORT"),
		RedisDB:                     getEnv("REDIS_DB"),
		RedisPassword:               getEnv("REDIS_PASSWORD"),*/
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}
