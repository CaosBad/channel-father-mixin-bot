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

type usersImpl struct{}

func registerUsers(router *httptreemux.TreeMux) {
	impl := &usersImpl{}
	router.POST("/auth", impl.authenticate)
	router.GET("/me", impl.me)
}

func (impl *usersImpl) authenticate(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
	} else if user, err := models.AuthenticateUserByOAuth(r.Context(), body.Code); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderUserWithAuthentication(w, r, user)
	}
}

func (impl *usersImpl) me(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderUserWithAuthentication(w, r, middlewares.CurrentUser(r))
}
