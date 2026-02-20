package llm

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// InitLogger configures logrus for Springfield.
// This is called from main() to ensure consistent logging across all packages.
func InitLogger() {
	// Configure logrus for our needs
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "15:04:05.000",
		FullTimestamp:   true,
	})

	// Set log level based on DEBUG environment variable
	if os.Getenv("DEBUG") == "1" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func init() {
	// Initialize logger on package import
	// This ensures logging is configured even if InitLogger() isn't called explicitly
	InitLogger()
}

// GetLogger returns a configured logrus entry with the given context name.
func GetLogger(context string) *log.Entry {
	return log.WithField("ctx", context)
}
