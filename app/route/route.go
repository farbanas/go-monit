package route

import (
	"github.com/julienschmidt/httprouter"
	"github.com/farbanas/go-monit/app/controllers"
	"github.com/farbanas/go-monit/app/controllers/handlerwrapper"
)

func Init() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", handlerwrapper.HandleFunc(controllers.Overview))
	router.POST("/webhooks/slack/monitor", handlerwrapper.HandleFunc(controllers.SlackMonitor))
	router.POST("/webhooks/load", handlerwrapper.HandleFunc(controllers.LoadSummary))
	router.POST("/webhooks/memory", handlerwrapper.HandleFunc(controllers.MemoryUsage))
	router.POST("/webhooks/process/:process", handlerwrapper.HandleFunc(controllers.ProcessStatus))

	return router
}