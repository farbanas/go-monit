package main

import (
	"log"
	"net/http"

	"github.com/farbanas/go-monit/app/utils"
	"github.com/farbanas/go-monit/app/route"
)

func main() {
	utils.Init()
	router := route.Init()

	log.Fatal(http.ListenAndServe(":10000", router))
}

