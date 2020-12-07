package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/wildangbudhi/forex-daily-analytics/depedencyinjection"
	"github.com/wildangbudhi/forex-daily-analytics/utils"
)

func main() {

	runtime.GOMAXPROCS(5)

	server, err := utils.NewServer()

	if err != nil {
		panic(err)
	}

	injectdepedency(server)

	// var wg sync.WaitGroup
	// wg.Add(2)

	// go func() {
	// 	defer wg.Done()
	// 	server.Router.Run(":8080")
	// }()

	// go func() {
	// defer wg.Done()
	server.Scheduler.Start()
	fmt.Println(server.Scheduler.Entries())
	// }()

	// wg.Wait()

	time.Sleep(5 * time.Minute)

}

func injectdepedency(server *utils.Server) {
	depedencyinjection.EconomicsData(server)
}
