package app

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateBot     CreateBotHandler
	DeleteBot     DeleteBotHandler
	StartBot      StartBotHandler
	StopBot       StopBotHandler
	UpdateStatus  UpdateStatusHandler
	Entry         EntryHandler
	Process       ProcessHandler
	CreateMailing CreateMailingHandler
	StartMailing  StartMailingHandler
}

type Queries struct {
	AllAnswers  GetAnswersTableHandler
	GetBot      GetBotHandler
	GetBots     GetBotsHandler
	StartedBots GetStartedBotsHandler
}
