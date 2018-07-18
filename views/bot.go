package views

import (
	"net/http"
	"time"

	"github.com/crossle/channel-father-mixin-bot/models"
)

type BotView struct {
	BotId      string    `json:"bot_id"`
	ExpireAt   time.Time `json:"expire_at"`
	ClientId   string    `json:"client_id"`
	SessionId  string    `json:"session_id"`
	PrivateKey string    `json:"private_key"`
}

func buildBotView(bot *models.Bot) BotView {
	view := BotView{
		BotId:      bot.BotId,
		ClientId:   bot.ClientId,
		SessionId:  bot.SessionId,
		PrivateKey: bot.PrivateKey,
		ExpireAt:   bot.ExpireAt,
	}
	return view
}

func RenderBot(w http.ResponseWriter, r *http.Request, bot *models.Bot) {
	RenderDataResponse(w, r, buildBotView(bot))
}
