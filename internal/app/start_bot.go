package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/pkg/decorator"
)

type StartBot struct {
	AuthorUUID string
	BotUUID    string
}

type StartBotHandler decorator.CommandHandler[StartBot]

type startBotHandler struct {
	bots bots.Repository
	im   bots.InstanceManager
}

func NewStartBotHandler(
	bots bots.Repository,
	im bots.InstanceManager,

	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) StartBotHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if im == nil {
		panic("instance manager is nil")
	}

	return decorator.ApplyCommandDecorators[StartBot](
		&startBotHandler{bots: bots, im: im},
		log,
		metricsClient,
	)
}

func (h startBotHandler) Handle(ctx context.Context, cmd StartBot) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	if err = bot.CanSeeBot(cmd.AuthorUUID); err != nil {
		return err
	}

	return h.im.Start(ctx, bot.UUID, bot.Token)
}
