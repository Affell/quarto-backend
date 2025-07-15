package config

import (
	"embed"
	"html/template"
	"os"
	"quarto/email"

	"github.com/charmbracelet/log"
	"github.com/provectio/godotenv"
)

const Version = "0.0.1"

var Config struct {
	FrontURL      string
	ListenPort    string
	BodySizeLimit string
	Email         email.Config
}

func Init(publicFolder embed.FS) {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found, getting variables from system environement")
	}

	frontURL := os.Getenv("FRONT_URL")
	if frontURL == "" {
		log.Fatal("FRONT_URL not set")
	}
	Config.FrontURL = frontURL

	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Warn("Invalid log level, setting to INFO")
		level = log.InfoLevel
	}
	log.SetLevel(level)

	listenPort, ok := os.LookupEnv("LISTEN_PORT")
	if !ok || listenPort == "" {
		log.Warn("LISTEN_PORT not set, using default value (5000)")
		listenPort = "5000"
	}
	Config.ListenPort = listenPort

	bodySizeLimit, ok := os.LookupEnv("MAX_BODY_SIZE")
	if !ok || bodySizeLimit == "" {
		log.Warn("MAX_BODY_SIZE not set or invalid, using default value (100 Mo)")
		bodySizeLimit = "100M"
	}
	Config.BodySizeLimit = bodySizeLimit

	if env := os.Getenv("SMTP_HOST"); env != "" {
		Config.Email.Host = env
	} else {
		log.Fatal("'SMTP_HOST' not in env")
	}

	if env := os.Getenv("SMTP_PORT"); env != "" {
		Config.Email.Port = env
	} else {
		log.Fatal("'SMTP_PORT' not in env")
	}

	if env := os.Getenv("SMTP_USER"); env != "" {
		Config.Email.User = env
	} else {
		log.Fatal("'SMTP_USER' not in env")
	}

	if env := os.Getenv("SMTP_PASSWORD"); env != "" {
		Config.Email.Password = env
	} else {
		log.Fatal("'SMTP_PASSWORD' not in env")
	}

	if env := os.Getenv("SMTP_DISPLAYNAME"); env != "" {
		Config.Email.From = env
	} else {
		log.Fatal("'SMTP_DISPLAYNAME' not in env")
	}

	Config.Email.RecoverTemplate = template.Must(template.ParseFS(publicFolder, "public/emails/recover.html"))
}
