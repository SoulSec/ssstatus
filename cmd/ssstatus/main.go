package main

import (
	"flag"
	"io/ioutil"
	"time"

	"github.com/SoulSec/ssstatus"
	"github.com/SoulSec/ssstatus/logger"
)

var configPath = flag.String("config", "configs/default.json", "configuration file")
var logPath = flag.String("log", "logs/from-"+time.Now().Format("2018-08-08")+".log", "log file")
var address = flag.String("http", ":8080", "address for http server")
var nolog = flag.Bool("nolog", false, "disable logging to file only")
var logfilter = flag.String("logfilter", "", "text to filter log by (both console and file)")

func main() {
	flag.Parse()
	jsonData, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic("error reading from configuration file")
	}

	if *nolog == true {
		logger.Disable()
	}

	if *logfilter != "" {
		logger.Filter(*logfilter)
	}

	logger.SetFilename(*logPath)

	config := ssstatus.NewConfig(jsonData)
	monitor := ssstatus.NewMonitor(config)
	go ssstatus.RunHttp(*address, monitor)
	monitor.Run()
}
