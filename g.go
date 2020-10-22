package main

import (
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
)

var LBConfig Config

var LBNodeP *LBNode
var LBTargetsMgrP *LBTargetsMgr
var LBConnectionPairMgrP *LBConnectionPairMgr

var LBGoroutineManagerP *goroutine_mgr.GoroutineManager
