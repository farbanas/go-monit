package main

import (
	"fmt"
	"sync"

	"git.vingd.com/v-lab/go-monit/machineinfo"
)

var wg = sync.WaitGroup{}

func main() {
	loadChan := make(chan []float64)

	go machineinfo.CoreLoad(loadChan)
	wg.Add(1)
	go func() {
		for {
			loads := <-loadChan
			for i, load := range loads {
				fmt.Printf("Core %d load: %.2f%%\n", i, load*100)
			}
		}
		wg.Done()
	}()

	wg.Wait()

	//http.HandleFunc("/", overview)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
func overview(w http.ResponseWriter, r *http.Request) {
	var loadChan chan float64 = make(chan float64, 5)
	machineinfo.CoreLoad(loadChan)
	return <-loadChan
}
*/
