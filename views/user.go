package views

import (
	"net/http"

	"github.com/crossle/channel-father-mixin-bot/models"
)

type UserView struct {
	Type           string `json:"type"`
	UserId         string `json:"user_id"`
	IdentityNumber string `json:"identity_number"`
	FullName       string `json:"full_name"`
	AvatarURL      string `json:"avatar_url"`
}

type UserWithAuthenticationView struct {
	UserView
	AuthenticationToken string `json:"authentication_token"`
}

func buildUserView(user *models.User) UserView {
	userView := UserView{
		Type:      "user",
		UserId:    user.UserId,
		FullName:  user.FullName,
		AvatarURL: user.AvatarURL,
	}
	return userView
}

func RenderUserWithAuthentication(w http.ResponseWriter, r *http.Request, user *models.User) {
	userView := UserWithAuthenticationView{
		UserView:            buildUserView(user),
		AuthenticationToken: user.AuthenticationToken,
	}
	RenderDataResponse(w, r, userView)
}
