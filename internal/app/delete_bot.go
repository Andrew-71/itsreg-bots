package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/pkg/decorator"
)

type DeleteBot struct {
	AuthorUUID string
	BotUUID    string
}

type DeleteBotHandler decorator.CommandHandler[DeleteBot]

type deleteBotHandler struct {
	bots bots.Repository
	im   bots.InstanceManager
}

func NewDeleteBotHandler(
	bots bots.Repository,
	im bots.InstanceManager,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) DeleteBotHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	return decorator.ApplyCommandDecorators[DeleteBot](
		deleteBotHandler{bots, im},
		logger, metricsClient,
	)
}

func (h deleteBotHandler) Handle(ctx context.Context, cmd DeleteBot) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	if err = bot.CanSeeBot(cmd.AuthorUUID); err != nil {
		return err
	}

	err = h.im.Stop(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	return h.bots.Delete(ctx, cmd.BotUUID)
}
