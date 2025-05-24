package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/pkg/decorator"
)

type StopBot struct {
	AuthorUUID string
	BotUUID    string
}

type StopBotHandler decorator.CommandHandler[StopBot]

type stopBotHandler struct {
	bots bots.Repository
	im   bots.InstanceManager
}

func NewStopBotHandler(
	bots bots.Repository,
	im bots.InstanceManager,

	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) StopBotHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if im == nil {
		panic("instance manager is nil")
	}

	return decorator.ApplyCommandDecorators[StopBot](
		&stopBotHandler{bots: bots, im: im},
		log,
		metricsClient,
	)
}

func (h stopBotHandler) Handle(ctx context.Context, cmd StopBot) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	if err = bot.CanSeeBot(cmd.AuthorUUID); err != nil {
		return err
	}

	return h.im.Stop(ctx, cmd.BotUUID)
}
