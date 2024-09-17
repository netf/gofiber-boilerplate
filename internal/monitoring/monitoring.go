package monitoring

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
)

func InitSentry(dsn string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
	})
	if err != nil {
		log.Error().Err(err).Msg("Sentry initialization failed")
	} else {
		log.Info().Msg("Sentry initialized")
	}
}

func FlushSentry() {
	sentry.Flush(2 * time.Second)
}
