package main

import (
	"fmt"
	"log"
	"net/http"

	"git.vingd.com/v-lab/go-monit/machineinfo"
)

func main() {
	http.HandleFunc("/", overview)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func overview(w http.ResponseWriter, r *http.Request) {
	var loadChan chan float64 = make(chan float64, 5)
	machineinfo.CoreLoad(loadChan)
	fmt.Fprintf(w, "%.2f", <-loadChan)
}
