package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/farbanas/go-monit/app/controllers"
	"github.com/farbanas/go-monit/app/utils/machineinfo"
)

func main() {
	loadChan := make(chan []float64)

	go coreLoadFeed(loadChan)
	go collectLoad(loadChan)

	http.HandleFunc("/", controllers.OverviewHandler)
	http.HandleFunc("/webhooks/slack/monitor", controllers.SlackMonitorHandler)
	http.HandleFunc("/webhooks/load", controllers.LoadSummaryHandler)
	http.HandleFunc("/webhooks/memory", controllers.MemoryUsageHandler)

	log.Fatal(http.ListenAndServe(":10000", nil))
}

func coreLoadFeed(loadChan chan []float64) {
	var prevCoreInfo []machineinfo.CoreInfo
	var curCoreInfo []machineinfo.CoreInfo

	curCoreInfo = machineinfo.ParseProcInfo()
	time.Sleep(1 * time.Second)
	for range time.Tick(1 * time.Second) {
		go func() {
			prevCoreInfo = curCoreInfo
			curCoreInfo = machineinfo.ParseProcInfo()
			loadChan <- machineinfo.CoreLoad(prevCoreInfo, curCoreInfo)
		}()
	}
}

func collectLoad(loadChan chan []float64) {
	for {
		loads := <-loadChan
		fmt.Println(loads)
	}
}
