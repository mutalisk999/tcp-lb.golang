package main

import (
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"runtime"
)

func initNode(c *Config) {
	LBNodeP = new(LBNode)
	LBNodeP.Initialise(c.Node.ListenEndPoint, c.Node.MaxConn, c.Node.Timeout)
}

func initTargetsMgr(c *Config) {
	LBTargetsMgrP = new(LBTargetsMgr)
	LBTargetsMgrP.Initialise()
	for _, t := range c.Targets {
		targetP := new(LBTarget)
		targetP.Initialise(t.ConnEndPoint, t.MaxConn, t.Timeout)

		targetId := CaclTargetId(t.ConnEndPoint)
		LBTargetsMgrP.Set(targetId, targetP)
	}
}

func initConnectionPairMgr() {
	LBConnectionPairMgrP = new(LBConnectionPairMgr)
	LBConnectionPairMgrP.Initialise()
}

func initGoroutineMgr() {
	LBGoroutineManagerP = new(goroutine_mgr.GoroutineManager)
	LBGoroutineManagerP.Initialise("global_goroutine_mgr")
}

func initApplication(c *Config) {
	initNode(c)
	initTargetsMgr(c)
	initConnectionPairMgr()
	initGoroutineMgr()

	if LBConfig.Threads > 0 {
		runtime.GOMAXPROCS(int(LBConfig.Threads))
		Info.Printf("Running with %v threads", LBConfig.Threads)
	}
}

func main() {
	iLogFile := "info.log"
	eLogFile := "error.log"
	InitLog(iLogFile, eLogFile, DEBUG)

	loadConfig(&LBConfig)
	initApplication(&LBConfig)

	LBGoroutineManagerP.GoroutineCreateP1("tcp_proxy_listener", startTcpProxy, &LBConfig)
	LBGoroutineManagerP.GoroutineCreateP1("api_server", startApiServer, &LBConfig)
	LBGoroutineManagerP.GoroutineCreateP0("maintain_loop", startMaintainLoop)

	quit := make(chan bool)
	<-quit
}
