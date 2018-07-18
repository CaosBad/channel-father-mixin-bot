package routes

import (
	"encoding/json"
	"net/http"

	"github.com/crossle/channel-father-mixin-bot/middlewares"
	"github.com/crossle/channel-father-mixin-bot/models"
	"github.com/crossle/channel-father-mixin-bot/session"
	"github.com/crossle/channel-father-mixin-bot/views"
	"github.com/dimfeld/httptreemux"
)

type botsImpl struct{}

func registerBots(router *httptreemux.TreeMux) {
	impl := &botsImpl{}
	router.POST("/bot/:id/keys", impl.postKeys)
	router.GET("/verify/:id", impl.verifyTrace)
	router.GET("/bot", impl.getBot)
}

func (impl *botsImpl) postKeys(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body struct {
		ClientId   string `json:"client_id"`
		SessionId  string `json:"session_id"`
		PrivateKey string `json:"private_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	current := middlewares.CurrentUser(r)
	botId := params["id"]
	if bot, err := models.PostKeys(r.Context(), current.UserId, botId, body.ClientId, body.SessionId, body.PrivateKey); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBot(w, r, bot)
	}
}

func (impl *botsImpl) getBot(w http.ResponseWriter, r *http.Request, params map[string]string) {
	current := middlewares.CurrentUser(r)
	if bot, err := models.FindBotByUserId(r.Context(), current.UserId); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if bot == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderBot(w, r, bot)
	}
}

func (impl *botsImpl) verifyTrace(w http.ResponseWriter, r *http.Request, params map[string]string) {
	current := middlewares.CurrentUser(r)
	if bot, err := models.VerifyTrace(r.Context(), current, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if bot == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderBot(w, r, bot)
	}
}
