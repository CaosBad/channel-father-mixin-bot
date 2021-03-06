package routes

import (
	"net/http"
	"runtime"
	"time"

	"github.com/crossle/channel-father-mixin-bot/config"
	"github.com/crossle/channel-father-mixin-bot/views"
	"github.com/dimfeld/httptreemux"
)

func RegisterRoutes(router *httptreemux.TreeMux) {
	router.GET("/", root)
	registerUsers(router)
	registerBots(router)
}
func root(w http.ResponseWriter, r *http.Request, params map[string]string) {
	now := time.Now()
	views.RenderDataResponse(w, r, map[string]string{
		"build":   config.BuildVersion + "-" + runtime.Version(),
		"website": "https://channel.otcxin.one",
		"iso":     now.UTC().Format(time.RFC3339Nano),
	})
}
