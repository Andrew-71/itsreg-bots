package service

import (
	"errors"
	"log/slog"
	"os"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"

	"github.com/bmstu-itstech/itsreg-bots/internal/app"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/infra"
	"github.com/bmstu-itstech/itsreg-bots/internal/service/mocks"
	"github.com/bmstu-itstech/itsreg-bots/pkg/decorator"
	"github.com/bmstu-itstech/itsreg-bots/pkg/logs"
	"github.com/bmstu-itstech/itsreg-bots/pkg/logs/handlers/slogdiscard"
	"github.com/bmstu-itstech/itsreg-bots/pkg/metrics"
)

func NewApplication() (
	app *app.Application,
	msgCh <-chan *message.Message,
	runCh <-chan *message.Message,
	close func() error,
) {
	logger := logs.DefaultLogger()
	metricsClient := metrics.NoOp{}

	url := os.Getenv("DATABASE_URI")
	db := sqlx.MustConnect("postgres", url)

	botsR := infra.NewPgBotsRepository(db)
	participants := infra.NewPgParticipantsRepository(db)

	msgPub, msgCh, senderClose := infra.NewNATSMessagesPublisher()
	runPub, runCh, senderClose := infra.NewNATSRunnerPublisher()

	return newApplication(
			logger, metricsClient, botsR, participants, msgPub, runPub,
		), msgCh, runCh, func() error {
			var err error
			err = errors.Join(err, db.Close())
			err = errors.Join(err, senderClose())
			return err
		}
}

func NewComponentTestApplication() (
	app *app.Application,
	msgCh <-chan *message.Message,
	runCh <-chan *message.Message,
) {
	logger := slogdiscard.NewDiscardLogger()
	metricsClient := metrics.NoOp{}

	botsR := mocks.NewMockBotRepository()
	participants := mocks.NewMockParticipantsRepository()

	msgPub, msgCh := mocks.NewMockMessagesPublisher()
	runPub, runCh := mocks.NewMockRunnerPublisher()

	return newApplication(
		logger, metricsClient, botsR, participants, msgPub, runPub,
	), msgCh, runCh
}

func newApplication(
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
	bots bots.Repository,
	participants bots.ParticipantRepository,
	msgPub bots.MessagesPublisher,
	runPub bots.RunnerPublisher,
) *app.Application {
	return &app.Application{
		Commands: app.Commands{
			CreateBot:     app.NewCreateBotHandler(bots, logger, metricsClient),
			DeleteBot:     app.NewDeleteBotHandler(bots, runPub, logger, metricsClient),
			StartBot:      app.NewStartBotHandler(bots, runPub, logger, metricsClient),
			StopBot:       app.NewStopBotHandler(bots, runPub, logger, metricsClient),
			UpdateStatus:  app.NewUpdateStatusHandler(bots, logger, metricsClient),
			Entry:         app.NewEntryHandler(bots, participants, msgPub, logger, metricsClient),
			Process:       app.NewProcessHandler(bots, participants, msgPub, logger, metricsClient),
			CreateMailing: app.NewCreateMailingHandler(bots, logger, metricsClient),
			StartMailing:  app.NewStartMailingHandler(bots, participants, msgPub, logger, metricsClient),
		},
		Queries: app.Queries{
			AllAnswers:  app.NewGetAnswersTableHandler(bots, participants, logger, metricsClient),
			GetBot:      app.NewGetBotHandler(bots, logger, metricsClient),
			GetBots:     app.NewGetBotsHandler(bots, logger, metricsClient),
			StartedBots: app.NewGetStartedBotsHandler(bots, logger, metricsClient),
		},
	}
}
