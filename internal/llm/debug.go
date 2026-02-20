package llm

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
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
		log.SetLevel(log.WarnLevel)
	}
}

// GetLogger returns a configured logrus entry with the given context name.
func GetLogger(context string) *log.Entry {
	return log.WithField("ctx", context)
}
