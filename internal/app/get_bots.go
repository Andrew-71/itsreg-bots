package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/pkg/decorator"
)

type GetBots struct {
	UserUUID string
}

type GetBotsHandler decorator.QueryHandler[GetBots, []Bot]

type getBotsHandler struct {
	bots bots.Repository
}

func NewGetBotsHandler(
	bots bots.Repository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetBotsHandler {
	return decorator.ApplyQueryDecorators[GetBots, []Bot](
		getBotsHandler{bots: bots},
		logger,
		metricsClient,
	)
}

func (h getBotsHandler) Handle(ctx context.Context, query GetBots) ([]Bot, error) {
	bs, err := h.bots.UserBots(ctx, query.UserUUID)
	if err != nil {
		return nil, err
	}

	return MapBotsFromDomain(bs), nil
}
