package route

import (
	"github.com/julienschmidt/httprouter"
	"github.com/farbanas/go-monit/app/controllers"
	"github.com/farbanas/go-monit/app/controllers/handlerwrapper"
)

func Init() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", handlerwrapper.HandleFunc(controllers.Overview))
	router.GET("/webhooks/slack/monitor", handlerwrapper.HandleFunc(controllers.SlackMonitor))
	router.GET("/webhooks/load", handlerwrapper.HandleFunc(controllers.LoadSummary))
	router.GET("/webhooks/memory", handlerwrapper.HandleFunc(controllers.MemoryUsage))
	router.GET("/webhooks/process/:process", handlerwrapper.HandleFunc(controllers.ProcessStatus))

	return router
}