package main

import (
	"fmt"
	"log"
	"net/http"

	"git.vingd.com/v-lab/go-monit/machineinfo"
)

func main() {
	machineinfo.ParseProcinfo()
	http.HandleFunc("/", overview)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func overview(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test")
}
