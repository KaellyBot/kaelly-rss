package application

import (
	"errors"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-rss/services/rss"
	"github.com/rs/zerolog/log"
)

var (
	ErrCannotInstanciateApp = errors.New("Cannot instanciate application")
)

type ApplicationInterface interface {
	Run() error
	Shutdown()
}

type Application struct {
	rss    rss.RSSServiceInterface
	broker amqp.MessageBrokerInterface
}

func New(rabbitMqClientId, rabbitMqAddress string, rssTimeout int) (*Application, error) {
	broker, err := amqp.New(rabbitMqClientId, rabbitMqAddress, []amqp.Binding{})
	if err != nil {
		log.Error().Err(err).Msgf("Failed to instanciate broker")
		return nil, ErrCannotInstanciateApp
	}

	rss, err := rss.New(broker, rssTimeout)
	if err != nil {
		log.Error().Err(err).Msgf("RSS service instanciation failed")
		return nil, ErrCannotInstanciateApp
	}

	return &Application{rss: rss, broker: broker}, nil
}

func (app *Application) Run() {
	app.rss.CheckFeeds()
}

func (app *Application) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
