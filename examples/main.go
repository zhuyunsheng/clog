package main

import (
	"github.com/clog"
	"strconv"
)

func main() {
	config := clog.Get("LOG")
	size, _ := strconv.Atoi(config["log_maxsize"])
	if config["log_file"] != "" {
		clog.InitLogger(config["log_file"], int64(size))
	} else {
		clog.InitLogger("log.log", int64(size))
	}
	logLevel, _ := strconv.Atoi(config["log_level"])
	clog.InitLogLevel(logLevel)
	clog.LogInfo("Info")
	clog.LogDebug("Debug")
	clog.LogErr("Error")
}
