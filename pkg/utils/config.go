package utils

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type smtpcreds struct {
	User     string `env:"USER"`
	Password string `env:"PASS"`
}

type cfg struct {
	Port             string      `env:"PORT" envDefault:"8080"`
	JwtSecret        string      `env:"JWT_SECRET,notEmpty"`
	PostgresHost     string      `env:"POSTGRES_HOST,notEmpty"`
	PostgresPort     string      `env:"POSTGRES_PORT,notEmpty"`
	PostgresUser     string      `env:"POSTGRES_USER,notEmpty"`
	PostgresPassword string      `env:"POSTGRES_PASSWORD,notEmpty"`
	PostgresDB       string      `env:"POSTGRES_DB,notEmpty"`
	RedisHost        string      `env:"REDIS_HOST,notEmpty"`
	RedisPort        string      `env:"REDIS_PORT,notEmpty"`
	RedisPassword    string      `env:"REDIS_PASSWORD,notEmpty"`
	EmailHost        string      `env:"EMAIL_HOST,notEmpty"`
	EmailPort        int         `env:"EMAIL_PORT,notEmpty"`
	SmtpCreds        []smtpcreds `envPrefix:"MAIL"`
	SendingEmail     string      `env:"SENDING_EMAIL,notEmpty"`
	RepoOwner        string      `env:"REPO_OWNER,notEmpty"`
	RepoName         string      `env:"REPO_NAME,notEmpty"`
	Recipients       string      `env:"RECIPIENETS"`
	CookieSecure     bool        `env:"SECURE" envDefault:"false"`
	GithubPAT        string      `env:"GITHUB_PAT"`
}

var Config cfg

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("No .env file found")
	}

	if err := env.Parse(&Config); err != nil {
		fmt.Printf("%+v", err)
		panic(err);
	}
}
