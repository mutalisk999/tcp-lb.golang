package main

import (
	"crypto/md5"
	"encoding/hex"
	"runtime"
)

func initNode(c *Config) {
	LBNodeP = new(LBNode)
	LBNodeP.Initialise(c.Node.ListenEndPoint, c.Node.MaxConn, c.Node.TimeoutRead)
}

func initTargetsMgr(c *Config) {
	LBTargetsMgrP = new(LBTargetsMgr)
	LBTargetsMgrP.Initialise()
	for _, t := range c.Targets {
		targetP := new(LBTarget)
		targetP.Initialise(t.ConnEndPoint, t.MaxConn, t.TimeoutRead, t.TimeoutRead)

		md5res := md5.Sum([]byte(t.ConnEndPoint))
		targetId := hex.EncodeToString(md5res[:])

		LBTargetsMgrP.Set(targetId, targetP)
	}
}

func initConnectionPairMgr() {
	LBConnectionPairMgrP = new(LBConnectionPairMgr)
	LBConnectionPairMgrP.Initialise()
}

func initApplication(c *Config) {
	initNode(c)
	initTargetsMgr(c)
	initConnectionPairMgr()
}

func main() {
	iLogFile := "info.log"
	eLogFile := "error.log"
	InitLog(iLogFile, eLogFile, DEBUG)
	LoadConfig(&LBConfig)

	if LBConfig.Threads > 0 {
		runtime.GOMAXPROCS(int(LBConfig.Threads))
		Info.Printf("Running with %v threads", LBConfig.Threads)
	}

	initApplication(&LBConfig)

	StartTcpProxy(&LBConfig)
	StartApiServer(&LBConfig)

	quit := make(chan bool)
	<-quit
}
