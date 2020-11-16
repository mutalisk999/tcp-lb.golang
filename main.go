package main

import (
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"runtime"
)

func initNode(c *Config) {
	LBNodeP = new(LBNode)
	if !verifyEndPointStr(c.Node.ListenEndPoint) {
		Error.Fatalf("invalid node listen endpoint: [%s]", c.Node.ListenEndPoint)
	}
	LBNodeP.Initialise(c.Node.ListenEndPoint, c.Node.MaxConn, c.Node.Timeout)
}

func initTargetsMgr(c *Config) {
	LBTargetsMgrP = new(LBTargetsMgr)
	LBTargetsMgrP.Initialise()
	for _, t := range c.Targets {
		targetP := new(LBTarget)
		if !verifyEndPointStr(t.ConnEndPoint) {
			Error.Fatalf("invalid target connect endpoint: [%s]", t.ConnEndPoint)
		}
		targetP.Initialise(t.ConnEndPoint, t.Active, t.MaxConn, t.Timeout)

		targetId := calcTargetId(t.ConnEndPoint)
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

	if !verifyEndPointStr(c.Api.ListenEndPoint) {
		Error.Fatalf("invalid api listen endpoint: [%s]", c.Api.ListenEndPoint)
	}

	if LBConfig.Threads > 0 {
		runtime.GOMAXPROCS(int(LBConfig.Threads))
		Info.Printf("Running with %v threads", LBConfig.Threads)
	}
}

func main() {
	loadConfig(&LBConfig)

	iLogFile := "info.log"
	eLogFile := "error.log"
	InitLog(iLogFile, eLogFile, LBConfig.Log.LogSetLevel)

	initApplication(&LBConfig)

	LBGoroutineManagerP.GoroutineCreateP1("tcp_proxy_listener", startTcpProxy, &LBConfig)
	LBGoroutineManagerP.GoroutineCreateP1("api_server", startApiServer, &LBConfig)
	LBGoroutineManagerP.GoroutineCreateP0("maintain_loop", startMaintainLoop)

	quit := make(chan bool)
	<-quit
}
