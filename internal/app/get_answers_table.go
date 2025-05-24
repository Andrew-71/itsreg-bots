package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/pkg/decorator"
)

type GetAnswersTable struct {
	UserUUID string
	BotUUID  string
}

type GetAnswersTableHandler decorator.QueryHandler[GetAnswersTable, AnswersTable]

type answersHandler struct {
	bots         bots.Repository
	participants bots.ParticipantRepository
}

func NewGetAnswersTableHandler(
	bots bots.Repository,
	participants bots.ParticipantRepository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetAnswersTableHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if participants == nil {
		panic("participants repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetAnswersTable, AnswersTable](
		answersHandler{bots: bots, participants: participants},
		logger,
		metricsClient,
	)
}

func (h answersHandler) Handle(ctx context.Context, query GetAnswersTable) (AnswersTable, error) {
	bot, err := h.bots.Bot(ctx, query.BotUUID)
	if err != nil {
		return AnswersTable{}, err
	}

	if err = bot.CanSeeBot(query.UserUUID); err != nil {
		return AnswersTable{}, err
	}

	prts, err := h.participants.ParticipantsOfBot(ctx, query.BotUUID)
	if err != nil {
		return AnswersTable{}, err
	}

	table := bots.NewAnswersTable(bot, prts)

	return MapAnswersTableFromDomain(table), nil
}
